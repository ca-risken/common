package dlp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const connectionYAMLTemplate = `
sources:
  fs:
    fs1:
      path: "%s"
      exclude_patterns:
        - "fingerprint.yaml"
`

// Scanner executes DLP scans on local directories
type Scanner struct {
	config *Config
}

// NewScanner creates a new Scanner instance with the given configuration
func NewScanner(config *Config) *Scanner {
	return &Scanner{
		config: config,
	}
}

// ScanDirectory scans the specified directory and returns the scan results
// dir: Local directory path containing files to scan
// resourceName: Resource name (e.g., bucket name, project name)
// totalFiles: Total number of files to be scanned
func (s *Scanner) ScanDirectory(ctx context.Context, dir string, resourceName string, totalFiles int) (*ScanResult, error) {
	startTime := time.Now()
	outputFile := filepath.Join(dir, "hawkeye-results.json")

	// Create connection file
	connectionFile := filepath.Join(dir, "connection.yml")
	connectionContent := fmt.Sprintf(connectionYAMLTemplate, dir)
	if err := os.WriteFile(connectionFile, []byte(connectionContent), 0600); err != nil {
		return nil, fmt.Errorf("failed to write connection file: %w", err)
	}

	// Copy fingerprint file
	fingerprintFile, err := s.config.CopyFingerprintFile(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to setup fingerprint file: %w", err)
	}

	// Execute hawk_scanner
	args := []string{"all", "--connection", connectionFile, "--fingerprint", fingerprintFile, "--shutup", "--json", outputFile}
	cmd := exec.CommandContext(ctx, "hawk_scanner", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "TERM=xterm") // Set TERM environment variable to avoid warning

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		outMsg := stdout.String()
		return nil, fmt.Errorf("hawk-eye scan failed: %w, stderr=%s, stdout=%s", err, errMsg, outMsg)
	}

	scanDuration := time.Since(startTime)
	return s.processScanResults(ctx, outputFile, resourceName, dir, totalFiles, scanDuration, startTime)
}

// processScanResults reads hawk-eye results and processes them
func (s *Scanner) processScanResults(ctx context.Context, outputFile, resourceName, tempDir string, totalFiles int, scanDuration time.Duration, scanTime time.Time) (*ScanResult, error) {
	// Check if results file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("hawk-eye results file not found: %s", outputFile)
	}

	// Read results file
	resultsData, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read results file: %w", err)
	}

	if len(resultsData) == 0 {
		return &ScanResult{
			ResourceName: resourceName,
			TotalFiles:   totalFiles,
			Findings:     []Finding{},
			ScanDuration: scanDuration.String(),
			ScanTime:     scanTime.Unix(),
		}, nil
	}

	// Parse hawk-eye results
	var hawkEyeOutput hawkEyeOutput
	if err := json.Unmarshal(resultsData, &hawkEyeOutput); err != nil {
		return nil, fmt.Errorf("failed to parse hawk-eye results as JSON: %w", err)
	}

	findings := []Finding{}
	for _, finding := range hawkEyeOutput.Fs {
		path := filepath.Clean(finding.FilePath)
		if tempDir != "" {
			path = filepath.Clean(finding.FilePath)
			path = filepath.Clean("/" + path[len(tempDir):])
		}

		// Get rule configuration
		rule := s.config.GetRule(finding.PatternName)
		if rule == nil {
			rule = &Rule{
				Name:        "Unknown",
				Description: "Unknown",
			}
		}

		// Check if rule is applicable to this file based on file filters
		fileInfo, err := os.Stat(finding.FilePath)
		if err == nil {
			if !rule.IsApplicableToFile(finding.FilePath, fileInfo.Size()) {
				continue
			}
		}

		// Limit matches to maximum matches per finding
		matches := finding.Matches
		totalMatchCount := len(matches)
		if len(matches) > s.config.MaxMatchesPerFinding {
			remainingCount := len(matches) - s.config.MaxMatchesPerFinding
			matches = matches[:s.config.MaxMatchesPerFinding]
			matches = append(matches, fmt.Sprintf("... and %d more", remainingCount))
		}

		findings = append(findings, Finding{
			FilePath:    filepath.Join(resourceName, path),
			Type:        rule.Type,
			PatternName: rule.Name,
			Description: rule.Description,
			Matches:     matches,
			Severity:    rule.CalculateSeverity(totalMatchCount),
		})
	}

	// Calculate total severity based on individual finding severities
	totalSeverity, severityReason := calculateTotalSeverity(findings)

	// Create structured scan result
	scanResult := &ScanResult{
		ResourceName:   resourceName,
		TotalFiles:     totalFiles,
		Findings:       findings,
		ScanDuration:   scanDuration.String(),
		ScanTime:       scanTime.Unix(),
		TotalSeverity:  totalSeverity,
		SeverityReason: severityReason,
	}
	return scanResult, nil
}

// calculateTotalSeverity determines the overall severity based on individual findings
func calculateTotalSeverity(findings []Finding) (string, string) {
	if len(findings) == 0 {
		return SeverityLow, "No findings detected"
	}

	// Count findings by severity
	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, finding := range findings {
		switch finding.Severity {
		case SeverityCritical:
			criticalCount++
		case SeverityHigh:
			highCount++
		case SeverityMedium:
			mediumCount++
		case SeverityLow:
			lowCount++
		}
	}

	// Create summary reason with all severity counts
	reason := fmt.Sprintf("CRITICAL: %d, HIGH: %d, MEDIUM: %d, LOW: %d", criticalCount, highCount, mediumCount, lowCount)

	// Determine total severity based on conditions
	switch {
	case criticalCount >= 1 || highCount >= 5:
		return SeverityCritical, reason
	case highCount >= 1 || mediumCount >= 10:
		return SeverityHigh, reason
	case mediumCount >= 1:
		return SeverityMedium, reason
	default: // LOW: All other cases
		return SeverityLow, reason
	}
}
