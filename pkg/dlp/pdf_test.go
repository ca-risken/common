package dlp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/ca-risken/core/proto/ai"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
)

type classifyPDFPageFunc func(context.Context, *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error)

func (f classifyPDFPageFunc) ClassifySensitiveContent(ctx context.Context, req *ai.ClassifySensitiveContentRequest, _ ...grpc.CallOption) (*ai.ClassifySensitiveContentResponse, error) {
	return f(ctx, req)
}

func TestPDFScanner(t *testing.T) {
	var inputs []string
	classifier := classifyPDFPageFunc(func(_ context.Context, req *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
		if req.GetFileName() != "test.pdf" {
			t.Errorf("file name = %q", req.GetFileName())
		}
		if req.GetText() != "" {
			inputs = append(inputs, "text")
			if req.GetText() != "password=example-only" || len(req.GetImage()) != 0 {
				t.Error("the text page must send the original extracted text")
			}
			return &ai.ClassifySensitiveContentResponse{
				Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
				Categories:  []ai.SensitiveContentCategory{ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_CREDENTIAL},
				ContentKind: ai.ContentKind_CONTENT_KIND_DOCUMENT,
				TestData:    ai.TestData_TEST_DATA_TRUE,
			}, nil
		}
		inputs = append(inputs, "image")
		config, err := jpeg.DecodeConfig(bytes.NewReader(req.GetImage()))
		if err != nil {
			t.Fatalf("decode JPEG: %v", err)
		}
		if max(config.Width, config.Height) != 2048 || len(req.GetImage()) > 2<<20 {
			t.Errorf("JPEG outside Core limits: %dx%d, %d bytes", config.Width, config.Height, len(req.GetImage()))
		}
		return &ai.ClassifySensitiveContentResponse{Sensitivity: ai.Sensitivity_SENSITIVITY_FALSE}, nil
	})
	scanner := NewPDFScanner(classifier)
	t.Cleanup(func() { closeTestPDFScanner(t, scanner) })

	got, err := scanner.Scan(t.Context(), "folder/test.pdf", bytes.NewReader(samplePDF()))
	if err != nil {
		t.Fatal(err)
	}
	want := PDFScan{
		Object: "folder/test.pdf", TotalPages: 2, Complete: true,
		Pages: []PDFPage{
			{Number: 1, Input: "text", Sensitivity: "TRUE", Categories: []string{"CREDENTIAL"}, ContentKind: "DOCUMENT", TestData: "TRUE"},
			{Number: 2, Input: "image", Sensitivity: "FALSE", TestData: "UNKNOWN"},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("PDF scan mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"text", "image"}, inputs); diff != "" {
		t.Errorf("Core calls mismatch (-want +got):\n%s", diff)
	}
}

func TestPDFScannerRechecksPublicAccessBeforeUsingCache(t *testing.T) {
	t.Parallel()

	updated := time.Date(2026, time.November, 20, 0, 0, 0, 0, time.UTC)
	previous := []PDFScan{
		{Object: "private.pdf", Updated: updated, Complete: true},
		{Object: "public.pdf", Updated: updated, Complete: true},
	}
	objects := []PDFObject{
		{
			Name: "private.pdf", Updated: updated,
			Open: func(context.Context) (io.ReadCloser, error) {
				return nil, errors.New("anonymous read denied")
			},
		},
		{
			Name: "public.pdf", Updated: updated,
			Open: func(context.Context) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("")), nil
			},
		},
	}
	scanner := &PDFScanner{}
	got := scanner.ScanObjects(t.Context(), objects, previous)
	want := []PDFScan{
		{Object: "private.pdf", Updated: updated},
		previous[1],
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("cached PDF scans mismatch (-want +got):\n%s", diff)
	}
}

func BenchmarkPDFScanner(b *testing.B) {
	input := samplePDF()
	classifier := classifyPDFPageFunc(func(context.Context, *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
		return &ai.ClassifySensitiveContentResponse{Sensitivity: ai.Sensitivity_SENSITIVITY_FALSE}, nil
	})
	b.Run("cold", func(b *testing.B) {
		for b.Loop() {
			scanner := NewPDFScanner(classifier)
			if _, err := scanner.Scan(b.Context(), "test.pdf", bytes.NewReader(input)); err != nil {
				b.Fatal(err)
			}
			closeTestPDFScanner(b, scanner)
		}
	})
	b.Run("warm", func(b *testing.B) {
		scanner := NewPDFScanner(classifier)
		defer closeTestPDFScanner(b, scanner)
		if _, err := scanner.Scan(b.Context(), "test.pdf", bytes.NewReader(input)); err != nil {
			b.Fatal(err)
		}
		for b.Loop() {
			if _, err := scanner.Scan(b.Context(), "test.pdf", bytes.NewReader(input)); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func closeTestPDFScanner(t testing.TB, scanner *PDFScanner) {
	t.Helper()
	pool, err := scanner.pool()
	if err != nil {
		t.Error(err)
		return
	}
	if err := pool.Close(); err != nil {
		t.Error(err)
	}
}

func samplePDF() []byte {
	stream := func(data string) string {
		return fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(data), data)
	}
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R 4 0 R] /Count 2 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 250 300] /Resources << /Font << /F1 5 0 R >> >> /Contents 6 0 R >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 250 300] /Resources << /XObject << /Im0 8 0 R >> >> /Contents 7 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		stream("BT /F1 12 Tf 10 200 Td (password=example-only) Tj ET"),
		stream("q 150 0 0 100 20 20 cm /Im0 Do Q"),
		"<< /Type /XObject /Subtype /Image /Width 1 /Height 1 /ColorSpace /DeviceRGB /BitsPerComponent 8 /Length 3 >>\nstream\n\xff\x00\x00\nendstream",
	}
	var pdf strings.Builder
	pdf.WriteString("%PDF-1.4\n")
	offsets := make([]int, 0, len(objects))
	for i, object := range objects {
		offsets = append(offsets, pdf.Len())
		fmt.Fprintf(&pdf, "%d 0 obj\n%s\nendobj\n", i+1, object)
	}
	xref := pdf.Len()
	fmt.Fprintf(&pdf, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for _, offset := range offsets {
		fmt.Fprintf(&pdf, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&pdf, "trailer\n<< /Root 1 0 R /Size %d >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xref)
	return []byte(pdf.String())
}
