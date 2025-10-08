package dlp

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

//go:embed yaml/dlp.yaml
var embeddedDLPYaml embed.FS

//go:embed yaml/fingerprint.yaml
var embeddedFingerprintYaml embed.FS

const (
	defaultDLPConfigFile  = "yaml/dlp.yaml"
	defaultFingerprintFile = "yaml/fingerprint.yaml"
)

// LoadConfig loads DLP configuration from file or default embedded configuration
// If configPath is empty, loads the default embedded configuration
func LoadConfig(configPath string) (*Config, error) {
	var yamlData []byte
	var err error

	if configPath != "" {
		// Load from external file
		configPath = filepath.Clean(configPath)
		yamlData, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read DLP config file %s: %w", configPath, err)
		}
	} else {
		// Load default embedded configuration
		yamlData, err = embeddedDLPYaml.ReadFile(defaultDLPConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded DLP config: %w", err)
		}
	}

	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse DLP config YAML: %w", err)
	}

	// Validate configuration using validator
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("invalid DLP configuration: %w", err)
	}
	return &config, nil
}

// GetMaxScanSizeBytes returns the maximum scan size in bytes
func (c *Config) GetMaxScanSizeBytes() int64 {
	if c.MaxScanSizeMB <= 0 {
		return 10 * 1024 * 1024 // Default 10MB if invalid
	}
	return int64(c.MaxScanSizeMB) * 1024 * 1024
}

// GetMaxSingleFileSizeBytes returns the maximum single file size in bytes
func (c *Config) GetMaxSingleFileSizeBytes() int64 {
	if c.MaxSingleFileSizeMB <= 0 {
		return 5 * 1024 * 1024 // Default 5MB if invalid
	}
	return int64(c.MaxSingleFileSizeMB) * 1024 * 1024
}

// GetFingerprintFilePath returns the fingerprint file path
// If not configured, returns the default embedded fingerprint file path
func (c *Config) GetFingerprintFilePath() string {
	if c.FingerprintFilePath != "" {
		return c.FingerprintFilePath
	}
	return defaultFingerprintFile
}

// CopyFingerprintFile copies the fingerprint file to the specified directory and returns the file path
func (c *Config) CopyFingerprintFile(destDir string) (string, error) {
	fingerprintConfigPath := c.GetFingerprintFilePath()
	if fingerprintConfigPath == "" {
		return "", fmt.Errorf("no fingerprint file configured")
	}
	destFile := filepath.Join(destDir, "fingerprint.yaml")

	// Read fingerprint data
	var fingerprintData []byte
	var err error

	if fingerprintConfigPath == defaultFingerprintFile {
		// Use embedded file
		fingerprintData, err = embeddedFingerprintYaml.ReadFile(fingerprintConfigPath)
		if err != nil {
			return "", fmt.Errorf("failed to read embedded fingerprint file: %w", err)
		}
	} else {
		// Use external file
		fingerprintConfigPath = filepath.Clean(fingerprintConfigPath)
		fingerprintData, err = os.ReadFile(fingerprintConfigPath)
		if err != nil {
			return "", fmt.Errorf("failed to read fingerprint file %s: %w", fingerprintConfigPath, err)
		}
	}

	// Write to destination
	if err := os.WriteFile(destFile, fingerprintData, 0600); err != nil {
		return "", fmt.Errorf("failed to write fingerprint file: %w", err)
	}
	return destFile, nil
}

// GetRule returns the rule for the given pattern name
func (c *Config) GetRule(patternName string) *Rule {
	for _, rule := range c.Rules {
		if rule.Name == patternName {
			return &rule
		}
	}
	return nil
}

// CalculateSeverity returns the severity level based on the match count
func (r *Rule) CalculateSeverity(matchCount int) string {
	// If no thresholds are configured, default to LOW
	if r.SeverityThresholds == nil {
		return SeverityLow
	}

	// Check thresholds from highest to lowest severity
	if r.SeverityThresholds.Critical != nil && matchCount >= *r.SeverityThresholds.Critical {
		return SeverityCritical
	}
	if r.SeverityThresholds.High != nil && matchCount >= *r.SeverityThresholds.High {
		return SeverityHigh
	}
	if r.SeverityThresholds.Medium != nil && matchCount >= *r.SeverityThresholds.Medium {
		return SeverityMedium
	}
	if r.SeverityThresholds.Low != nil && matchCount >= *r.SeverityThresholds.Low {
		return SeverityLow
	}
	// Default to LOW if no threshold matches
	return SeverityLow
}

// IsApplicableToFile checks if this rule should be applied to the given file
func (r *Rule) IsApplicableToFile(fileName string, fileSizeBytes int64) bool {
	// No filters means apply to all files
	if r.FileFilters == nil {
		return true
	}

	// Check minimum file size
	if r.FileFilters.MinSizeKB != nil {
		fileSizeKB := int(fileSizeBytes / 1024)
		if fileSizeKB < *r.FileFilters.MinSizeKB {
			return false
		}
	}

	// Check include extensions (if specified, file must match one of them)
	if len(r.FileFilters.IncludeExtensions) > 0 {
		ext := filepath.Ext(fileName)
		if !slices.Contains(r.FileFilters.IncludeExtensions, ext) {
			return false
		}
	}

	// Check exclude file name patterns
	if len(r.FileFilters.ExcludeFileName) > 0 {
		baseName := filepath.Base(fileName)
		for _, pattern := range r.FileFilters.ExcludeFileName {
			matched, err := filepath.Match(pattern, baseName)
			if err == nil && matched {
				return false
			}
		}
	}

	return true
}
