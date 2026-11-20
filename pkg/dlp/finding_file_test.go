package dlp

import (
	json "encoding/json/v2"
	"testing"

	"github.com/ca-risken/core/proto/ai"
	"github.com/google/go-cmp/cmp"
)

func TestNewFindingFile(t *testing.T) {
	t.Parallel()

	pii := ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_PII
	clientData := ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_CLIENT_DATA
	sensitive := ai.Sensitivity_SENSITIVITY_TRUE
	notSensitive := ai.Sensitivity_SENSITIVITY_FALSE
	testData := ai.TestData_TEST_DATA_TRUE

	tests := []struct {
		name string
		scan ScanResponse
		want FindingFile
	}{
		{
			name: "CSV with unknown test data",
			scan: ScanResponse{Complete: true, Classification: &Classification{
				Input: "text", Sensitivity: sensitive, ContentKind: ai.ContentKind_CONTENT_KIND_TABULAR_DATA,
				Categories: []ai.SensitiveContentCategory{pii, clientData},
			}},
			want: FindingFile{
				Input: "text", Complete: true, Sensitivity: new(true),
				ContentType: []string{"CLIENT_DATA", "PII"}, ContentKind: "TABULAR_DATA",
			},
		},
		{
			name: "PDF is one file with categories from all pages",
			scan: ScanResponse{Complete: true, PDFResult: &PDFResult{
				TotalPages: 3,
				Pages: []PDFPage{
					{Number: 1, Input: "text", Sensitivity: sensitive, Categories: []ai.SensitiveContentCategory{clientData}, TestData: testData},
					{Number: 2, Input: "image", Sensitivity: sensitive, Categories: []ai.SensitiveContentCategory{pii, clientData}, TestData: testData},
					{Number: 3, Input: "text", Sensitivity: notSensitive},
				},
			}},
			want: FindingFile{
				Input: "pdf", Complete: true, Sensitivity: new(true), TotalPages: 3,
				ContentType: []string{"CLIENT_DATA", "PII"}, ContentKind: "DOCUMENT", TestData: new(true),
			},
		},
		{
			name: "known non-sensitive file has a false decision",
			scan: ScanResponse{Complete: true, Classification: &Classification{Input: "text", Sensitivity: notSensitive}},
			want: FindingFile{
				Input: "text", Complete: true, Sensitivity: new(false),
				ContentType: []string{}, ContentKind: "UNSPECIFIED",
			},
		},
		{
			name: "known production data stores false",
			scan: ScanResponse{Complete: true, Classification: &Classification{
				Input: "text", Sensitivity: sensitive, TestData: ai.TestData_TEST_DATA_FALSE,
				Categories: []ai.SensitiveContentCategory{pii},
			}},
			want: FindingFile{
				Input: "text", Complete: true, Sensitivity: new(true),
				ContentType: []string{"PII"}, ContentKind: "UNSPECIFIED", TestData: new(false),
			},
		},
		{
			name: "partial negative is unknown for the file",
			scan: ScanResponse{Classification: &Classification{Input: "text", Truncated: true, Sensitivity: notSensitive}},
			want: FindingFile{
				Input: "text", Truncated: true,
				ContentType: []string{}, ContentKind: "UNSPECIFIED",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := NewFindingFile(test.scan)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("metadata mismatch (-want +got):\n%s", diff)
			}

			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]any
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			if _, present := fields["sensitivity"]; present {
				t.Error("file metadata includes sensitivity")
			}
			value, present := fields["test_data"]
			if test.want.TestData == nil {
				if present {
					t.Errorf("unknown test_data should be omitted, got %v", value)
				}
			} else if value != *test.want.TestData {
				t.Errorf("test_data = %v, want JSON boolean %t", value, *test.want.TestData)
			}

			wantSummary := NewFindingScan([]ScanResponse{test.scan}, nil)
			gotSummary := NewFindingScan([]ScanResponse{got.ScanResponse()}, nil)
			if diff := cmp.Diff(wantSummary, gotSummary); diff != "" {
				t.Errorf("restored summary mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
