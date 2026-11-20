package dlp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"slices"
	"strings"

	"github.com/ca-risken/core/proto/ai"
)

type Scanner struct {
	client ai.SensitiveContentServiceClient
	pdf    *pdfScanner
	config *Config
}

func NewScanner(client ai.SensitiveContentServiceClient, config *Config) *Scanner {
	return &Scanner{client: client, pdf: newPDFScanner(client), config: config}
}

// CanScanFile checks the configured file name and size limits before reading the file.
func (s *Scanner) CanScanFile(name string, size int64) bool {
	return size <= s.config.GetMaxSingleFileSizeBytes() && !slices.ContainsFunc(s.config.ExcludeFilePatterns, func(pattern string) bool {
		return matchesExcludePattern(name, pattern)
	})
}

// Scan classifies one file. It returns nil for an excluded, empty, or unsupported
// file. The caller owns the reader and any associated cleanup.
func (s *Scanner) Scan(ctx context.Context, req ScanRequest) (*ScanResponse, error) {
	if req.Reader == nil {
		return nil, fmt.Errorf("file reader is required")
	}
	if !s.CanScanFile(req.FileName, 0) {
		return nil, nil
	}

	maxBytes := s.config.GetMaxSingleFileSizeBytes()
	data, err := io.ReadAll(io.LimitReader(req.Reader, maxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("file exceeds the %d-byte limit", maxBytes)
	}
	if len(data) == 0 {
		return nil, nil
	}
	if strings.EqualFold(path.Ext(req.FileName), ".zip") {
		return s.scanZIP(ctx, data)
	}
	return s.scanContents(ctx, req.FileName, http.DetectContentType(data), data)
}

func (s *Scanner) scanContents(ctx context.Context, name, contentType string, data []byte) (*ScanResponse, error) {
	fileName := limitUTF8(path.Base(name), 1024)
	if contentType == "application/pdf" || strings.EqualFold(path.Ext(name), ".pdf") {
		return s.pdf.scan(ctx, fileName, data)
	}

	classification := &Classification{}
	coreReq := &ai.ClassifySensitiveContentRequest{FileName: fileName}
	switch {
	case contentType == "image/jpeg" || contentType == "image/png":
		classification.Input = "image"
		if len(data) > maxImageBytes {
			return &ScanResponse{Classification: classification}, nil
		}
		coreReq.Image = data
	case strings.HasPrefix(contentType, "text/") || hasTextExtension(name):
		text, ok := decodeText(name, data)
		if !ok {
			return nil, nil
		}
		classification.Input = "text"
		classification.Truncated = len(text) > maxTextBytes
		coreReq.Text = limitUTF8(text, maxTextBytes)
	default:
		return nil, nil
	}

	reviewed := classify(ctx, s.client, coreReq, classification)
	return &ScanResponse{Complete: reviewed && !classification.Truncated, Classification: classification}, nil
}

func matchesExcludePattern(name, pattern string) bool {
	if !strings.HasPrefix(pattern, ".") {
		return strings.Contains(name, pattern)
	}
	for part := range strings.SplitSeq(name, "/") {
		if strings.EqualFold(part, pattern) || strings.EqualFold(path.Ext(part), pattern) {
			return true
		}
	}
	return false
}

func classify(ctx context.Context, client ai.SensitiveContentServiceClient, req *ai.ClassifySensitiveContentRequest, result *Classification) bool {
	response, err := client.ClassifySensitiveContent(ctx, req)
	if err != nil || response == nil {
		return false
	}
	result.Sensitivity = response.GetSensitivity()
	result.Categories = response.GetCategories()
	result.ContentKind = response.GetContentKind()
	result.TestData = response.GetTestData()
	return true
}
