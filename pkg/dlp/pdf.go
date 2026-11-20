package dlp

import (
	"context"
	"errors"
	"io"
	"path"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/ca-risken/core/proto/ai"
	pdfium "github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
)

const (
	maxPDFBytes    = 5 << 20
	maxPDFPages    = 20
	maxTextBytes   = 32 << 10
	maxImageBytes  = 2 << 20
	maxImageEdge   = 2048
	pdfTimeout     = 2 * time.Minute
	pdfMemoryPages = 2048 // WebAssembly pages of 64 KiB each.
)

type PDFScan struct {
	Object     string    `json:"object"`
	Updated    time.Time `json:"updated,omitzero"`
	TotalPages int       `json:"total_pages,omitempty"`
	Pages      []PDFPage `json:"pages,omitempty"`
	Complete   bool      `json:"complete"`
}

type PDFObject struct {
	Name    string
	Updated time.Time
	Open    func(context.Context) (io.ReadCloser, error)
}

// PDFPage keeps the Core enum types. The cached JSON uses their numeric values.
type PDFPage struct {
	Number      int                           `json:"number"`
	Input       string                        `json:"input,omitempty"`
	Truncated   bool                          `json:"truncated,omitempty"`
	Sensitivity ai.Sensitivity                `json:"sensitivity"`
	Categories  []ai.SensitiveContentCategory `json:"categories,omitempty"`
	ContentKind ai.ContentKind                `json:"content_kind,omitzero"`
	TestData    ai.TestData                   `json:"test_data,omitzero"`
}

type PDFScanner struct {
	client ai.SensitiveContentServiceClient
	pool   func() (pdfium.Pool, error)
}

func NewPDFScanner(client ai.SensitiveContentServiceClient) *PDFScanner {
	return &PDFScanner{
		client: client,
		pool: sync.OnceValues(func() (pdfium.Pool, error) {
			return webassembly.Init(webassembly.Config{
				MaxTotal:     1,
				MaxIdle:      1,
				ReuseWorkers: true,
				FSConfig:     wazero.NewFSConfig(),
				Stdout:       io.Discard,
				Stderr:       io.Discard,
				RuntimeConfig: wazero.NewRuntimeConfig().
					WithCoreFeatures(api.CoreFeaturesV2 | experimental.CoreFeaturesExceptionHandling).
					WithMemoryLimitPages(pdfMemoryPages).
					WithCloseOnContextDone(true),
			})
		}),
	}
}

// IsPDFObject accepts PDF names from the bucket listing without downloading objects.
func IsPDFObject(name string) bool {
	return strings.EqualFold(path.Ext(name), ".pdf")
}

// ScanObjects reviews a bounded batch and reuses completed results for unchanged PDFs.
func (s *PDFScanner) ScanObjects(ctx context.Context, objects []PDFObject, previous []PDFScan) []PDFScan {
	ctx, cancel := context.WithTimeout(ctx, pdfTimeout)
	defer cancel()
	cached := make(map[string]PDFScan, len(previous))
	for _, scan := range previous {
		if scan.Complete && !scan.Updated.IsZero() {
			cached[scan.Object] = scan
		}
	}
	var results []PDFScan
	for _, object := range objects {
		result := PDFScan{Object: object.Name, Updated: object.Updated}
		if ctx.Err() == nil {
			reader, err := object.Open(ctx)
			if err == nil {
				if scan, ok := cached[object.Name]; ok && object.Updated.Equal(scan.Updated) {
					result = scan
				} else {
					result, _ = s.Scan(ctx, object.Name, reader)
				}
				_ = reader.Close()
				result.Updated = object.Updated
			}
		}
		results = append(results, result)
	}
	return results
}

// Scan reads a publicly readable PDF and records classification metadata for each
// attempted page. The caller retains ownership of the reader.
func (s *PDFScanner) Scan(ctx context.Context, object string, reader io.Reader) (PDFScan, error) {
	result := PDFScan{Object: object}
	ctx, cancel := context.WithTimeout(ctx, pdfTimeout)
	defer cancel()
	data, err := io.ReadAll(io.LimitReader(reader, maxPDFBytes+1))
	if err != nil || len(data) > maxPDFBytes || len(data) == 0 {
		return result, errors.New("could not read PDF within size limit")
	}
	pool, err := s.pool()
	if err != nil {
		return result, errors.New("could not initialize PDFium")
	}
	instance, err := pool.GetInstanceWithContext(ctx)
	if err != nil {
		return result, errors.New("could not acquire PDFium worker")
	}
	interruptDone := make(chan struct{})
	stopInterrupt := context.AfterFunc(ctx, func() {
		_ = instance.Kill()
		close(interruptDone)
	})
	defer func() {
		if stopInterrupt() {
			_ = instance.Close()
		} else {
			<-interruptDone
		}
	}()
	document, err := instance.OpenDocument(&requests.OpenDocument{File: &data})
	if err != nil {
		return result, errors.New("could not open PDF")
	}
	pages, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{Document: document.Document})
	if err != nil || pages.PageCount < 1 {
		return result, errors.New("could not count PDF pages")
	}
	result.TotalPages = pages.PageCount
	result.Complete = pages.PageCount <= maxPDFPages
	for i := range min(pages.PageCount, maxPDFPages) {
		if ctx.Err() != nil {
			result.Complete = false
			return result, ctx.Err()
		}
		page := requests.Page{ByIndex: new(requests.PageByIndex{Document: document.Document, Index: i})}
		classification, reviewed := s.scanPage(ctx, instance, page, limitUTF8(path.Base(object), 1024))
		classification.Number = i + 1
		result.Pages = append(result.Pages, classification)
		result.Complete = result.Complete && reviewed && !classification.Truncated
	}
	return result, nil
}

func (s *PDFScanner) scanPage(ctx context.Context, instance pdfium.Pdfium, page requests.Page, fileName string) (PDFPage, bool) {
	var result PDFPage
	req := &ai.ClassifySensitiveContentRequest{FileName: fileName}
	text, textErr := instance.GetPageText(&requests.GetPageText{Page: page})
	if textErr == nil && strings.TrimSpace(text.Text) != "" {
		result.Input = "text"
		result.Truncated = len(text.Text) > maxTextBytes
		req.Text = limitUTF8(text.Text, maxTextBytes)
	} else {
		result.Input = "image"
		rendered, err := instance.RenderToFile(&requests.RenderToFile{
			RenderPageInPixels: new(requests.RenderPageInPixels{Page: page, Width: maxImageEdge, Height: maxImageEdge}),
			OutputFormat:       requests.RenderToFileOutputFormatJPG,
			OutputTarget:       requests.RenderToFileOutputTargetBytes,
			OutputQuality:      85,
			MaxFileSize:        maxImageBytes,
		})
		if err != nil || rendered.ImageBytes == nil || len(*rendered.ImageBytes) == 0 || len(*rendered.ImageBytes) > maxImageBytes {
			return result, false
		}
		req.Image = *rendered.ImageBytes
	}
	response, err := s.client.ClassifySensitiveContent(ctx, req)
	if err != nil || response == nil {
		return result, false
	}
	result.Sensitivity = response.GetSensitivity()
	result.Categories = response.GetCategories()
	result.ContentKind = response.GetContentKind()
	result.TestData = response.GetTestData()
	return result, true
}

func limitUTF8(text string, maximum int) string {
	if len(text) <= maximum {
		return text
	}
	for !utf8.RuneStart(text[maximum]) {
		maximum--
	}
	return text[:maximum]
}
