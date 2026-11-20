package dlp

import (
	"maps"
	"slices"
	"strings"

	"github.com/ca-risken/core/proto/ai"
)

// FindingScan is the shared dlp_scan JSON payload. It includes the scan
// details needed to reuse results, plus a summary for Finding consumers.
type FindingScan struct {
	*ScanResult
	Sensitive   *bool    `json:"sensitive"`
	ContentType []string `json:"content_type"`
	TestData    *bool    `json:"test_data"`
}

// NewFindingScan summarizes known sensitive matches and PDF page decisions.
// Unknown or incomplete evidence leaves the corresponding flag null.
func NewFindingScan(scan *ScanResult) *FindingScan {
	if scan == nil {
		return nil
	}

	result := &FindingScan{ScanResult: scan}
	contentTypes := make(map[string]struct{})
	hasSensitive := false
	complete := len(scan.PDFScans) > 0
	for _, finding := range scan.Findings {
		if contentType := legacyContentType(finding.Type); contentType != "" {
			hasSensitive = true
			contentTypes[contentType] = struct{}{}
		} else {
			complete = false
		}
	}

	allSensitiveAreTest := !hasSensitive
	hasNonTestData := false
	for _, pdf := range scan.PDFScans {
		complete = complete && pdf.Complete && len(pdf.Pages) > 0
		for _, page := range pdf.Pages {
			switch page.Sensitivity {
			case ai.Sensitivity_SENSITIVITY_TRUE:
				hasSensitive = true
				hasNonTestData = hasNonTestData || page.TestData == ai.TestData_TEST_DATA_FALSE
				allSensitiveAreTest = allSensitiveAreTest && page.TestData == ai.TestData_TEST_DATA_TRUE
				for _, category := range page.Categories {
					if category == ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_UNSPECIFIED {
						continue
					}
					if name, ok := strings.CutPrefix(category.String(), "SENSITIVE_CONTENT_CATEGORY_"); ok {
						contentTypes[name] = struct{}{}
					}
				}
			case ai.Sensitivity_SENSITIVITY_FALSE:
			default:
				complete = false
			}
		}
	}

	if hasSensitive || complete {
		result.Sensitive = new(hasSensitive)
	}
	switch {
	case hasNonTestData:
		result.TestData = new(false)
	case hasSensitive && complete && allSensitiveAreTest:
		result.TestData = new(true)
	}
	result.ContentType = slices.AppendSeq(make([]string, 0, len(contentTypes)), maps.Keys(contentTypes))
	slices.Sort(result.ContentType)
	return result
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
