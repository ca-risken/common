package dlp

import (
	"io"

	"github.com/ca-risken/core/proto/ai"
)

// ScanRequest represents the name and contents of one file.
type ScanRequest struct {
	FileName string
	Reader   io.Reader
}

// ScanResponse holds the classification of a file or its individual PDF pages.
type ScanResponse struct {
	Complete bool `json:"complete"`
	*Classification
	*PDFResult
}

type Classification struct {
	Input       string                        `json:"input,omitempty"`
	Truncated   bool                          `json:"truncated,omitempty"`
	Sensitivity ai.Sensitivity                `json:"sensitivity"`
	Categories  []ai.SensitiveContentCategory `json:"categories,omitempty"`
	ContentKind ai.ContentKind                `json:"content_kind,omitzero"`
	TestData    ai.TestData                   `json:"test_data,omitzero"`
}

type PDFResult struct {
	TotalPages int       `json:"total_pages,omitempty"`
	Pages      []PDFPage `json:"pages,omitempty"`
}

type PDFPage struct {
	Number int `json:"number"`
	Classification
}

// Finding represents a previously stored directory scanner match.
type Finding struct {
	FilePath    string   `json:"file_path"`
	Type        string   `json:"type"`
	PatternName string   `json:"pattern_name"`
	Matches     []string `json:"matches"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
}

type Config struct {
	MaxSingleFileSizeMB int      `yaml:"max_single_file_size_mb" validate:"required,gt=0"`
	ExcludeFilePatterns []string `yaml:"exclude_file_patterns"`
}
