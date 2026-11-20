package dlp

import (
	"bytes"
	"compress/gzip"
	"context"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/ai"
	"github.com/google/go-cmp/cmp"
)

func TestScanGZIPLog(t *testing.T) {
	t.Parallel()

	const log = "2026-11-21T01:02:03Z user=person@example.test action=login"
	calls := 0
	classifier := classifyPDFPageFunc(func(_ context.Context, req *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
		calls++
		if req.GetFileName() != "access.LOG" || req.GetText() != log || len(req.GetImage()) != 0 {
			t.Error("expected Core to receive the text and object basename without the gzip extension")
		}
		return &ai.ClassifySensitiveContentResponse{
			Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
			Categories:  []ai.SensitiveContentCategory{ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_PII},
			ContentKind: ai.ContentKind_CONTENT_KIND_LOG,
			TestData:    ai.TestData_TEST_DATA_TRUE,
		}, nil
	})

	scan, err := NewScanner(classifier, &Config{}).Scan(t.Context(), ScanRequest{
		FileName: "logs/access.LOG.GZ",
		Reader:   bytes.NewReader(makeTestGZIP(t, []byte(log))),
	})
	if err != nil || scan == nil {
		t.Fatalf("scan gzip log: response=%v, error=%v", scan, err)
	}
	if calls != 1 {
		t.Errorf("Core calls=%d, want 1", calls)
	}
	want := FindingFile{
		Input:       "text",
		Complete:    true,
		Sensitivity: new(true),
		ContentType: []string{"PII"},
		ContentKind: "LOG",
		TestData:    new(true),
	}
	if diff := cmp.Diff(want, NewFindingFile(*scan)); diff != "" {
		t.Errorf("log Finding differs (-want +got):\n%s", diff)
	}
}

func TestScanGZIPLogEligibilityAndLimits(t *testing.T) {
	t.Parallel()

	compressed := makeTestGZIP(t, []byte("sample log"))
	for _, test := range []struct {
		name      string
		fileName  string
		input     []byte
		wantError string
	}{
		{name: "other gzip file", fileName: "archive.tar.gz", input: compressed},
		{name: "binary disguised as a log", fileName: "access.log.gz", input: makeTestGZIP(t, []byte("%PDF-1.7\n"))},
		{name: "truncated gzip log", fileName: "access.log.gz", input: compressed[:len(compressed)-1], wantError: "read gzip log"},
		{
			name:      "expanded log exceeds limit",
			fileName:  "access.log.gz",
			input:     makeTestGZIP(t, bytes.Repeat([]byte{'x'}, (1<<20)+1)),
			wantError: "expanded gzip log exceeds",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			classifier := classifyPDFPageFunc(func(context.Context, *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
				t.Fatal("file should be rejected before calling Core")
				return nil, nil
			})
			scan, err := NewScanner(classifier, &Config{MaxSingleFileSizeMB: 1}).Scan(t.Context(), ScanRequest{
				FileName: test.fileName,
				Reader:   bytes.NewReader(test.input),
			})
			if scan != nil {
				t.Errorf("scan=%v, want nil", scan)
			}
			if test.wantError == "" {
				if err != nil {
					t.Errorf("scan error=%v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("scan error=%v, want %q", err, test.wantError)
			}
		})
	}
}

func makeTestGZIP(t *testing.T, data []byte) []byte {
	t.Helper()

	var output bytes.Buffer
	writer := gzip.NewWriter(&output)
	writer.Name = "header.log"
	if _, err := writer.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}
