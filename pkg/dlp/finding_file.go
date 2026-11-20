package dlp

import (
	"slices"
	"strings"

	"github.com/ca-risken/core/proto/ai"
)

// FindingFile is the file-level classification written to Finding JSON.
type FindingFile struct {
	Input       string   `json:"input,omitempty"`
	Complete    bool     `json:"complete"`
	Truncated   bool     `json:"truncated,omitempty"`
	Sensitivity *bool    `json:"-"`
	ContentType []string `json:"content_type"`
	ContentKind string   `json:"content_kind"`
	TestData    *bool    `json:"test_data,omitzero"`
	TotalPages  int      `json:"total_pages,omitempty"`
}

// NewFindingFile summarizes a scan as one file, including all of its reviewed PDF pages.
func NewFindingFile(scan ScanResponse) FindingFile {
	file := FindingFile{Complete: scan.Complete, ContentKind: "UNSPECIFIED"}
	file.applySummary(NewFindingScan([]ScanResponse{scan}, nil))
	if scan.PDFResult != nil {
		file.Input = "pdf"
		file.ContentKind = "DOCUMENT"
		file.TotalPages = scan.TotalPages
		file.Truncated = scan.TotalPages > maxPDFPages || slices.ContainsFunc(scan.Pages, func(page PDFPage) bool {
			return page.Truncated
		})
		return file
	}
	if scan.Classification != nil {
		file.Input = scan.Input
		file.Truncated = scan.Truncated
		if kind, ok := strings.CutPrefix(scan.ContentKind.String(), "CONTENT_KIND_"); ok {
			file.ContentKind = kind
		}
	}
	return file
}

// ScanResponse restores the file-level decisions for Finding aggregation.
func (f FindingFile) ScanResponse() ScanResponse {
	classification := &Classification{
		Input:       f.Input,
		Truncated:   f.Truncated,
		Sensitivity: ai.Sensitivity(ai.Sensitivity_value["SENSITIVITY_"+findingDecisionName(f.Sensitivity)]),
		ContentKind: ai.ContentKind(ai.ContentKind_value["CONTENT_KIND_"+f.ContentKind]),
		TestData:    ai.TestData(ai.TestData_value["TEST_DATA_"+findingDecisionName(f.TestData)]),
	}
	for _, name := range f.ContentType {
		value := ai.SensitiveContentCategory_value["SENSITIVE_CONTENT_CATEGORY_"+name]
		if value != 0 {
			classification.Categories = append(classification.Categories, ai.SensitiveContentCategory(value))
		}
	}
	return ScanResponse{Complete: f.Complete, Classification: classification}
}

// MergeLegacyFindings folds previous matches into the file summary without including their values.
func (f *FindingFile) MergeLegacyFindings(previous []Finding) {
	f.applySummary(NewFindingScan([]ScanResponse{f.ScanResponse()}, previous))
}

func (f *FindingFile) applySummary(summary FindingScan) {
	f.Sensitivity = summary.Sensitive
	f.ContentType = summary.ContentType
	f.TestData = summary.TestData
}

func findingDecisionName(value *bool) string {
	switch {
	case value == nil:
		return "UNKNOWN"
	case *value:
		return "TRUE"
	default:
		return "FALSE"
	}
}
