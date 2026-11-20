package dlp

import (
	json "encoding/json/v2"
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
		name string
		scan ScanResult
		want FindingScan
	}{
		{
			name: "all sensitive pages are test data",
			scan: ScanResult{PDFScans: []PDFScan{{Complete: true, Pages: []PDFPage{testPage, negativePage}}}},
			want: FindingScan{Sensitive: new(true), ContentType: []string{"CLIENT_DATA", "PII"}, TestData: new(true)},
		},
		{
			name: "a non-test page makes the summary non-test",
			scan: ScanResult{PDFScans: []PDFScan{{Complete: true, Pages: []PDFPage{testPage, nonTestPage}}}},
			want: FindingScan{Sensitive: new(true), ContentType: []string{"CLIENT_DATA", "CREDENTIAL", "PII"}, TestData: new(false)},
		},
		{
			name: "all reviewed pages are non-sensitive",
			scan: ScanResult{PDFScans: []PDFScan{{Complete: true, Pages: []PDFPage{negativePage}}}},
			want: FindingScan{Sensitive: new(false), ContentType: []string{}},
		},
		{
			name: "incomplete scan cannot declare data non-sensitive",
			scan: ScanResult{PDFScans: []PDFScan{{Pages: []PDFPage{negativePage}}}},
			want: FindingScan{ContentType: []string{}},
		},
		{
			name: "unknown page prevents the all-test conclusion",
			scan: ScanResult{PDFScans: []PDFScan{{Complete: true, Pages: []PDFPage{testPage, {}}}}},
			want: FindingScan{Sensitive: new(true), ContentType: []string{"CLIENT_DATA", "PII"}},
		},
		{
			name: "legacy detections contribute categories without a test-data decision",
			scan: ScanResult{Findings: []Finding{{Type: "Credential"}, {Type: "Webhook URL"}, {Type: "PII"}}},
			want: FindingScan{Sensitive: new(true), ContentType: []string{"CREDENTIAL", "PII"}},
		},
		{
			name: "a cloud identifier alone does not determine sensitivity",
			scan: ScanResult{Findings: []Finding{{Type: "Cloud Resource"}}},
			want: FindingScan{ContentType: []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			want := tt.want
			want.ScanResult = &tt.scan
			if diff := cmp.Diff(&want, NewFindingScan(&tt.scan)); diff != "" {
				t.Errorf("Finding scan mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFindingScanJSONPreservesCache(t *testing.T) {
	t.Parallel()

	if got := NewFindingScan(nil); got != nil {
		t.Fatal("a missing scan must stay unset")
	}
	scan := &ScanResult{
		ResourceName: "public-bucket",
		ScanTime:     100,
		Findings:     []Finding{{FilePath: "public-bucket/record", Type: "PII", Matches: []string{"person@example.com"}}},
		PDFScans: []PDFScan{{
			Object:   "document.pdf",
			Complete: true,
			Pages:    []PDFPage{{Sensitivity: ai.Sensitivity_SENSITIVITY_FALSE}},
		}},
	}
	data, err := json.Marshal(NewFindingScan(scan))
	if err != nil {
		t.Fatal(err)
	}
	var cached ScanResult
	if err := json.Unmarshal(data, &cached); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(scan, &cached); diff != "" {
		t.Errorf("cached scan mismatch (-want +got):\n%s", diff)
	}
}
