package dlp

import (
	"testing"

	"github.com/ca-risken/core/proto/ai"
	"github.com/google/go-cmp/cmp"
)

func TestNewFindingScan(t *testing.T) {
	t.Parallel()

	testPage := PDFPage{
		Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
		Categories: []ai.SensitiveContentCategory{
			ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_PII,
			ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_CLIENT_DATA,
			ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_PII,
			ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_UNSPECIFIED,
		},
		TestData: ai.TestData_TEST_DATA_TRUE,
	}
	nonTestPage := PDFPage{
		Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
		Categories:  []ai.SensitiveContentCategory{ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_CREDENTIAL},
		TestData:    ai.TestData_TEST_DATA_FALSE,
	}
	negativePage := PDFPage{Sensitivity: ai.Sensitivity_SENSITIVITY_FALSE}

	tests := []struct {
		name     string
		files    []ScanResponse
		previous []Finding
		want     FindingScan
	}{
		{
			name:  "all sensitive pages are test data",
			files: pdfResponses(true, testPage, negativePage),
			want:  FindingScan{Sensitive: new(true), ContentType: []string{"CLIENT_DATA", "PII"}, TestData: new(true)},
		},
		{
			name:  "a non-test page makes the summary non-test",
			files: pdfResponses(true, testPage, nonTestPage),
			want:  FindingScan{Sensitive: new(true), ContentType: []string{"CLIENT_DATA", "CREDENTIAL", "PII"}, TestData: new(false)},
		},
		{
			name:  "all reviewed pages are non-sensitive",
			files: pdfResponses(true, negativePage),
			want:  FindingScan{Sensitive: new(false), ContentType: []string{}},
		},
		{
			name:  "incomplete scan cannot declare data non-sensitive",
			files: pdfResponses(false, negativePage),
			want:  FindingScan{ContentType: []string{}},
		},
		{
			name:  "incomplete scan retains sensitive evidence",
			files: pdfResponses(false, testPage),
			want:  FindingScan{Sensitive: new(true), ContentType: []string{"CLIENT_DATA", "PII"}},
		},
		{
			name:  "unknown page prevents the all-test conclusion",
			files: pdfResponses(true, testPage, PDFPage{}),
			want:  FindingScan{Sensitive: new(true), ContentType: []string{"CLIENT_DATA", "PII"}},
		},
		{
			name:     "legacy detections contribute categories without a test-data decision",
			previous: []Finding{{Type: "Credential"}, {Type: "Webhook URL"}, {Type: "PII"}},
			want:     FindingScan{Sensitive: new(true), ContentType: []string{"CREDENTIAL", "PII"}},
		},
		{
			name:     "an unrecognized finding does not determine sensitivity",
			previous: []Finding{{Type: "Identifier"}},
			want:     FindingScan{ContentType: []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if diff := cmp.Diff(tt.want, NewFindingScan(tt.files, tt.previous)); diff != "" {
				t.Errorf("Finding scan mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func pdfResponses(complete bool, pages ...PDFPage) []ScanResponse {
	return []ScanResponse{{Complete: complete, PDFResult: &PDFResult{Pages: pages}}}
}
