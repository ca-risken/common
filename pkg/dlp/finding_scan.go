package dlp

import (
	"maps"
	"slices"
	"strings"

	"github.com/ca-risken/core/proto/ai"
)

// FindingScan is the shared classification summary for Finding consumers.
type FindingScan struct {
	Sensitive   *bool    `json:"sensitive"`
	ContentType []string `json:"content_type"`
	TestData    *bool    `json:"test_data"`
}

// NewFindingScan converts file decisions and previously stored matches to a
// Finding summary. Unknown or incomplete evidence leaves the corresponding flag null.
func NewFindingScan(files []ScanResponse, previousFindings []Finding) FindingScan {
	summary := findingScanSummary{
		contentTypes:        make(map[string]struct{}),
		complete:            len(files) > 0,
		allSensitiveAreTest: true,
	}
	for _, finding := range previousFindings {
		summary.addLegacyFinding(finding)
	}
	for _, file := range files {
		summary.addFile(file)
	}

	var result FindingScan
	if summary.hasSensitive || summary.complete {
		result.Sensitive = new(summary.hasSensitive)
	}
	switch {
	case summary.hasNonTestData:
		result.TestData = new(false)
	case summary.hasSensitive && summary.complete && summary.allSensitiveAreTest:
		result.TestData = new(true)
	}
	result.ContentType = slices.AppendSeq(make([]string, 0, len(summary.contentTypes)), maps.Keys(summary.contentTypes))
	slices.Sort(result.ContentType)
	return result
}

type findingScanSummary struct {
	contentTypes        map[string]struct{}
	complete            bool
	hasSensitive        bool
	allSensitiveAreTest bool
	hasNonTestData      bool
}

func (s *findingScanSummary) addLegacyFinding(finding Finding) {
	contentType := legacyContentType(finding.Type)
	if contentType == "" {
		s.complete = false
		return
	}

	s.hasSensitive = true
	s.allSensitiveAreTest = false
	s.contentTypes[contentType] = struct{}{}
}

func (s *findingScanSummary) addFile(file ScanResponse) {
	s.complete = s.complete && file.Complete
	if file.ZIP != nil {
		for _, entry := range file.ZIP.Entries {
			s.addFile(entry.Scan)
		}
		return
	}
	if file.PDFResult != nil {
		s.complete = s.complete && len(file.Pages) > 0
		for _, page := range file.Pages {
			s.addClassification(page.Classification)
		}
		return
	}
	if file.Classification == nil {
		s.complete = false
		return
	}
	s.addClassification(*file.Classification)
}

func (s *findingScanSummary) addClassification(result Classification) {
	if result.Sensitivity == ai.Sensitivity_SENSITIVITY_FALSE {
		return
	}
	if result.Sensitivity != ai.Sensitivity_SENSITIVITY_TRUE {
		s.complete = false
		return
	}

	s.hasSensitive = true
	s.hasNonTestData = s.hasNonTestData || result.TestData == ai.TestData_TEST_DATA_FALSE
	s.allSensitiveAreTest = s.allSensitiveAreTest && result.TestData == ai.TestData_TEST_DATA_TRUE
	s.addCategories(result.Categories)
}

func (s *findingScanSummary) addCategories(categories []ai.SensitiveContentCategory) {
	for _, category := range categories {
		if category == ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_UNSPECIFIED {
			continue
		}
		name, ok := strings.CutPrefix(category.String(), "SENSITIVE_CONTENT_CATEGORY_")
		if !ok {
			continue
		}
		s.contentTypes[name] = struct{}{}
	}
}

func legacyContentType(findingType string) string {
	switch findingType {
	case "PII":
		return "PII"
	case "Credential", "Webhook URL":
		return "CREDENTIAL"
	case "Client Info":
		return "CLIENT_DATA"
	default:
		return ""
	}
}
