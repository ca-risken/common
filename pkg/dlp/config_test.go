package dlp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		wantErr    bool
	}{
		{
			name:       "load embedded config",
			configPath: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadConfig(tt.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && config == nil {
				t.Error("LoadConfig() returned nil config")
			}
		})
	}
}

func TestCalculateSeverity(t *testing.T) {
	critical := 1
	high := 5
	medium := 10
	low := 1

	tests := []struct {
		name       string
		rule       Rule
		matchCount int
		want       string
	}{
		{
			name: "critical threshold",
			rule: Rule{
				Name: "Test Rule",
				SeverityThresholds: &SeverityThresholds{
					Critical: &critical,
					High:     &high,
					Medium:   &medium,
					Low:      &low,
				},
			},
			matchCount: 1,
			want:       SeverityCritical,
		},
		{
			name: "high threshold",
			rule: Rule{
				Name: "Test Rule",
				SeverityThresholds: &SeverityThresholds{
					High:   &high,
					Medium: &medium,
					Low:    &low,
				},
			},
			matchCount: 5,
			want:       SeverityHigh,
		},
		{
			name: "medium threshold",
			rule: Rule{
				Name: "Test Rule",
				SeverityThresholds: &SeverityThresholds{
					Medium: &medium,
					Low:    &low,
				},
			},
			matchCount: 10,
			want:       SeverityMedium,
		},
		{
			name: "low threshold",
			rule: Rule{
				Name: "Test Rule",
				SeverityThresholds: &SeverityThresholds{
					Low: &low,
				},
			},
			matchCount: 1,
			want:       SeverityLow,
		},
		{
			name: "no thresholds",
			rule: Rule{
				Name: "Test Rule",
			},
			matchCount: 100,
			want:       SeverityLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.CalculateSeverity(tt.matchCount)
			if got != tt.want {
				t.Errorf("CalculateSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsApplicableToFile(t *testing.T) {
	minSize := 1

	tests := []struct {
		name          string
		rule          Rule
		fileName      string
		fileSizeBytes int64
		want          bool
	}{
		{
			name: "no filters - always applicable",
			rule: Rule{
				Name: "Test Rule",
			},
			fileName:      "test.txt",
			fileSizeBytes: 1000,
			want:          true,
		},
		{
			name: "include extension - matches",
			rule: Rule{
				Name: "Test Rule",
				FileFilters: &FileFilters{
					IncludeExtensions: []string{".csv", ".txt"},
				},
			},
			fileName:      "test.txt",
			fileSizeBytes: 1000,
			want:          true,
		},
		{
			name: "include extension - not matches",
			rule: Rule{
				Name: "Test Rule",
				FileFilters: &FileFilters{
					IncludeExtensions: []string{".csv", ".log"},
				},
			},
			fileName:      "test.txt",
			fileSizeBytes: 1000,
			want:          false,
		},
		{
			name: "min size - above threshold",
			rule: Rule{
				Name: "Test Rule",
				FileFilters: &FileFilters{
					MinSizeKB: &minSize,
				},
			},
			fileName:      "test.txt",
			fileSizeBytes: 2048, // 2KB
			want:          true,
		},
		{
			name: "min size - below threshold",
			rule: Rule{
				Name: "Test Rule",
				FileFilters: &FileFilters{
					MinSizeKB: &minSize,
				},
			},
			fileName:      "test.txt",
			fileSizeBytes: 512, // 0.5KB
			want:          false,
		},
		{
			name: "exclude file name - matches pattern",
			rule: Rule{
				Name: "Test Rule",
				FileFilters: &FileFilters{
					ExcludeFileName: []string{"*test*"},
				},
			},
			fileName:      "test_file.txt",
			fileSizeBytes: 1000,
			want:          false,
		},
		{
			name: "exclude file name - not matches pattern",
			rule: Rule{
				Name: "Test Rule",
				FileFilters: &FileFilters{
					ExcludeFileName: []string{"*backup*"},
				},
			},
			fileName:      "data.txt",
			fileSizeBytes: 1000,
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.IsApplicableToFile(tt.fileName, tt.fileSizeBytes)
			if got != tt.want {
				t.Errorf("IsApplicableToFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateTotalSeverity(t *testing.T) {
	tests := []struct {
		name               string
		findings           []Finding
		wantSeverity       string
		wantReasonContains string
	}{
		{
			name:               "no findings",
			findings:           []Finding{},
			wantSeverity:       SeverityLow,
			wantReasonContains: "No findings detected",
		},
		{
			name: "1 critical finding",
			findings: []Finding{
				{Severity: SeverityCritical},
			},
			wantSeverity:       SeverityCritical,
			wantReasonContains: "CRITICAL: 1",
		},
		{
			name: "5 high findings",
			findings: []Finding{
				{Severity: SeverityHigh},
				{Severity: SeverityHigh},
				{Severity: SeverityHigh},
				{Severity: SeverityHigh},
				{Severity: SeverityHigh},
			},
			wantSeverity:       SeverityCritical,
			wantReasonContains: "HIGH: 5",
		},
		{
			name: "1 high finding",
			findings: []Finding{
				{Severity: SeverityHigh},
			},
			wantSeverity:       SeverityHigh,
			wantReasonContains: "HIGH: 1",
		},
		{
			name: "10 medium findings",
			findings: []Finding{
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
				{Severity: SeverityMedium},
			},
			wantSeverity:       SeverityHigh,
			wantReasonContains: "MEDIUM: 10",
		},
		{
			name: "1 medium finding",
			findings: []Finding{
				{Severity: SeverityMedium},
			},
			wantSeverity:       SeverityMedium,
			wantReasonContains: "MEDIUM: 1",
		},
		{
			name: "only low findings",
			findings: []Finding{
				{Severity: SeverityLow},
				{Severity: SeverityLow},
			},
			wantSeverity:       SeverityLow,
			wantReasonContains: "LOW: 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeverity, gotReason := calculateTotalSeverity(tt.findings)
			if gotSeverity != tt.wantSeverity {
				t.Errorf("calculateTotalSeverity() severity = %v, want %v", gotSeverity, tt.wantSeverity)
			}
			if diff := cmp.Diff(tt.wantReasonContains, gotReason); diff != "" && tt.wantReasonContains != gotReason {
				// Check if wantReasonContains is a substring of gotReason
				if len(tt.wantReasonContains) > 0 && !contains(gotReason, tt.wantReasonContains) {
					t.Errorf("calculateTotalSeverity() reason does not contain %v, got %v", tt.wantReasonContains, gotReason)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
