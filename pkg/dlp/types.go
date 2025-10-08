package dlp

// Severity level constants
const (
	SeverityLow      = "LOW"
	SeverityMedium   = "MEDIUM"
	SeverityHigh     = "HIGH"
	SeverityCritical = "CRITICAL"
)

// ScanResult represents the complete DLP scan result
type ScanResult struct {
	ResourceName   string    `json:"resource_name"`
	TotalFiles     int       `json:"total_files"`
	Findings       []Finding `json:"findings"`
	ScanDuration   string    `json:"scan_duration"`
	ScanTime       int64     `json:"scan_time"` // Unix timestamp
	TotalSeverity  string    `json:"total_severity"`
	SeverityReason string    `json:"severity_reason"`
}

// Finding represents an individual DLP finding
type Finding struct {
	FilePath    string   `json:"file_path"`
	Type        string   `json:"type"`
	PatternName string   `json:"pattern_name"`
	Matches     []string `json:"matches"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
}

// Config represents the complete DLP configuration
type Config struct {
	MaxScanFiles         int                 `yaml:"max_scan_files" validate:"required,gt=0"`
	MaxScanSizeMB        int                 `yaml:"max_scan_size_mb" validate:"required,gt=0"`
	MaxSingleFileSizeMB  int                 `yaml:"max_single_file_size_mb" validate:"required,gt=0"`
	MaxMatchesPerFinding int                 `yaml:"max_matches_per_finding" validate:"required,gt=0"`
	ExcludeFilePatterns  []string            `yaml:"exclude_file_patterns"`
	FingerprintFilePath  string              `yaml:"fingerprint_file_path,omitempty"`
	Rules                []Rule              `yaml:"rules,omitempty" validate:"dive"`
}

// Rule represents a single DLP rule
type Rule struct {
	Name               string              `yaml:"name" validate:"required,min=1"`
	Description        string              `yaml:"description"`
	Type               string              `yaml:"type"`
	FileFilters        *FileFilters        `yaml:"file_filters,omitempty"`
	SeverityThresholds *SeverityThresholds `yaml:"severity_thresholds,omitempty"`
}

// FileFilters defines file metadata conditions for applying a rule
type FileFilters struct {
	IncludeExtensions []string `yaml:"include_extensions,omitempty"`
	ExcludeFileName   []string `yaml:"exclude_file_name,omitempty"`
	MinSizeKB         *int     `yaml:"min_size_kb,omitempty"`
}

// SeverityThresholds defines match count thresholds for different severity levels
type SeverityThresholds struct {
	Critical *int `yaml:"critical,omitempty"`
	High     *int `yaml:"high,omitempty"`
	Medium   *int `yaml:"medium,omitempty"`
	Low      *int `yaml:"low,omitempty"`
}

// hawkEyeOutput represents the complete hawk-eye JSON output structure (internal use only)
type hawkEyeOutput struct {
	Fs []Finding `json:"fs"`
}
