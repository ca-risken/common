package dlp

import (
	"archive/zip"
	"bytes"
	"context"
	json "encoding/json/v2"
	"fmt"
	"image"
	"image/png"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/ai"
	"github.com/google/go-cmp/cmp"
)

func TestScanZIPAsSingleFindingFile(t *testing.T) {
	var pngData bytes.Buffer
	if err := png.Encode(&pngData, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	input := []zipTestEntry{
		{name: "reports/customers.csv", data: []byte("name,email\nExample,person@example.test\n")},
		{name: "reports/agreement.pdf", data: samplePDF()},
		{name: "images/card.png", data: pngData.Bytes()},
		{name: "notes.txt", data: []byte("Public notice")},
	}
	var calls []string
	classifier := classifyPDFPageFunc(func(_ context.Context, req *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
		calls = append(calls, req.FileName)
		if req.FileName == "card.png" && !bytes.Equal(req.Image, pngData.Bytes()) {
			t.Error("expected PNG bytes to reach Core")
		}
		if req.FileName == "notes.txt" || (req.FileName == "agreement.pdf" && req.Text == "") {
			return &ai.ClassifySensitiveContentResponse{Sensitivity: ai.Sensitivity_SENSITIVITY_FALSE}, nil
		}
		category := ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_PII
		if req.FileName == "agreement.pdf" {
			category = ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_CLIENT_DATA
		}
		return &ai.ClassifySensitiveContentResponse{
			Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
			Categories:  []ai.SensitiveContentCategory{category},
			TestData:    ai.TestData_TEST_DATA_TRUE,
		}, nil
	})
	scanner := NewScanner(classifier, &Config{})
	t.Cleanup(func() { closeTestPDFScanner(t, scanner.pdf) })
	scan, err := scanner.Scan(t.Context(), ScanRequest{
		FileName: "reports.ZIP",
		Reader:   bytes.NewReader(makeTestZIP(t, input)),
	})
	if err != nil || scan == nil {
		t.Fatalf("scan ZIP: response=%v, error=%v", scan, err)
	}
	want := FindingFile{
		Input:            "zip",
		Complete:         true,
		Sensitivity:      new(true),
		ContentType:      []string{"CLIENT_DATA", "PII"},
		ContentKind:      "OTHER",
		TestData:         new(true),
		SensitiveEntries: []string{"reports/customers.csv", "reports/agreement.pdf", "images/card.png"},
	}
	got := NewFindingFile(*scan)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("archive Finding differs (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"customers.csv", "agreement.pdf", "agreement.pdf", "card.png", "notes.txt"}, calls); diff != "" {
		t.Errorf("Core requests differ (-want +got):\n%s", diff)
	}
	wire, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(wire), "person@example.test") || strings.Contains(string(wire), `"pages"`) || strings.Contains(string(wire), `"notes.txt"`) {
		t.Error("Finding exposed contents or non-sensitive members")
	}
	var restored FindingFile
	if err := json.Unmarshal(wire, &restored); err != nil {
		t.Fatal(err)
	}
	restored.Sensitivity = new(true)
	if diff := cmp.Diff(want, restored); diff != "" {
		t.Errorf("stored archive differs (-want +got):\n%s", diff)
	}
}

func TestScanZIPEligibilityAndLimits(t *testing.T) {
	t.Parallel()

	text := []byte("synthetic value")
	nested := makeTestZIP(t, []zipTestEntry{{name: "inside.txt", data: text}})
	entries := make([]zipTestEntry, maxZIPEntries+1)
	for i := range entries {
		entries[i] = zipTestEntry{name: fmt.Sprintf("%d.txt", i), data: text}
	}
	chunk := bytes.Repeat([]byte{'x'}, 5<<20)
	big := make([]zipTestEntry, maxZIPExpandedBytes/len(chunk)+1)
	for i := range big {
		big[i] = zipTestEntry{name: fmt.Sprintf("%d.txt", i), data: chunk}
	}

	for _, test := range []struct {
		name          string
		input         []zipTestEntry
		wantCalls     int
		wantSkipped   bool
		wantTruncated bool
	}{
		{
			name:      "skip nested archives even when renamed and continue with text",
			input:     []zipTestEntry{{name: "inside.zip", data: nested}, {name: "inside.csv", data: nested}, {name: "tail.txt", data: text}},
			wantCalls: 1,
		},
		{
			name:        "skip entire archive when an entry is encrypted",
			input:       []zipTestEntry{{name: "plain.txt", data: text}, {name: "encrypted.txt", data: text, flags: 1}},
			wantSkipped: true,
		},
		{name: "entry count limit", input: entries, wantCalls: maxZIPEntries, wantTruncated: true},
		{name: "expanded byte limit", input: big, wantCalls: maxZIPExpandedBytes / len(chunk), wantTruncated: true},
		{
			name:          "expanded member exceeds file size limit",
			input:         []zipTestEntry{{name: "big.txt", data: append(bytes.Clone(chunk), 'x')}, {name: "tail.txt", data: text}},
			wantCalls:     1,
			wantTruncated: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			calls := 0
			classifier := classifyPDFPageFunc(func(context.Context, *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
				calls++
				return &ai.ClassifySensitiveContentResponse{
					Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
					TestData:    ai.TestData_TEST_DATA_TRUE,
				}, nil
			})
			scan, err := NewScanner(classifier, &Config{}).Scan(t.Context(), ScanRequest{
				FileName: "archive.zip",
				Reader:   bytes.NewReader(makeTestZIP(t, test.input)),
			})
			if err != nil || (scan == nil) != test.wantSkipped {
				t.Fatalf("scan=%v err=%v, want skipped=%t", scan, err, test.wantSkipped)
			}
			if calls != test.wantCalls {
				t.Errorf("Core calls=%d, want %d", calls, test.wantCalls)
			}
			if scan == nil {
				return
			}
			if scan.Complete || scan.ZIP.Truncated != test.wantTruncated {
				t.Errorf("complete=%t truncated=%t, want incomplete with truncated=%t", scan.Complete, scan.ZIP.Truncated, test.wantTruncated)
			}
			if finding := NewFindingFile(*scan); finding.Sensitivity == nil || !*finding.Sensitivity || finding.TestData != nil {
				t.Error("partial archive should retain detections with unknown test-data status")
			}
		})
	}
}

type zipTestEntry struct {
	name  string
	data  []byte
	flags uint16
}

func makeTestZIP(t *testing.T, entries []zipTestEntry) []byte {
	t.Helper()
	var output bytes.Buffer
	writer := zip.NewWriter(&output)
	for _, entry := range entries {
		member, err := writer.CreateHeader(&zip.FileHeader{Name: entry.name, Method: zip.Deflate, Flags: entry.flags})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := member.Write(entry.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}
