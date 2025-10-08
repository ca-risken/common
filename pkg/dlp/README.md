# DLP (Data Loss Prevention) Package

Cloud-agnostic DLP scanning package for detecting sensitive information in files using hawk-eye scanner.

## Overview

This package provides a common DLP scanning functionality that can be used across different cloud providers (AWS, GCP, Azure, etc.). It scans local directories for sensitive data such as credentials, PII, and other confidential information.

## Features

- **Cloud-agnostic**: Works with any cloud provider by scanning local files
- **Configurable rules**: Customizable detection patterns and severity thresholds
- **File filtering**: Include/exclude files based on extensions, size, and patterns
- **Severity calculation**: Automatic severity assessment based on match counts
- **Embedded defaults**: Built-in configuration and fingerprint files

## Installation

```bash
go get github.com/ca-risken/common/pkg/dlp
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ca-risken/common/pkg/dlp"
)

func main() {
    // Load DLP configuration (uses embedded default if path is empty)
    config, err := dlp.LoadConfig("")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Create scanner
    scanner := dlp.NewScanner(config)

    // Scan directory
    result, err := scanner.ScanDirectory(
        context.Background(),
        "/path/to/scan",
        "my-resource-name",
        100, // total files
    )
    if err != nil {
        log.Fatalf("Scan failed: %v", err)
    }

    // Process results
    fmt.Printf("Total Severity: %s\n", result.TotalSeverity)
    fmt.Printf("Findings: %d\n", len(result.Findings))
    for _, finding := range result.Findings {
        fmt.Printf("  - %s: %s (%s)\n",
            finding.FilePath,
            finding.PatternName,
            finding.Severity,
        )
    }
}
```

### AWS S3 Integration Example

```go
package main

import (
    "context"
    "io"
    "os"

    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/ca-risken/common/pkg/dlp"
)

func scanS3Bucket(ctx context.Context, s3Client *s3.Client, bucketName string) (*dlp.ScanResult, error) {
    // 1. Create temporary directory
    tempDir, err := os.MkdirTemp("", "dlp-scan-")
    if err != nil {
        return nil, err
    }
    defer os.RemoveAll(tempDir)

    // 2. Download files from S3 to local directory
    // (Your S3 download logic here)
    // ...

    // 3. Load DLP configuration
    config, err := dlp.LoadConfig("")
    if err != nil {
        return nil, err
    }

    // 4. Scan the directory
    scanner := dlp.NewScanner(config)
    return scanner.ScanDirectory(ctx, tempDir, bucketName, fileCount)
}
```

### GCP Cloud Storage Integration Example

```go
package main

import (
    "context"
    "io"
    "os"

    "cloud.google.com/go/storage"
    "github.com/ca-risken/common/pkg/dlp"
)

func scanGCSBucket(ctx context.Context, storageClient *storage.Client, bucketName string) (*dlp.ScanResult, error) {
    // 1. Create temporary directory
    tempDir, err := os.MkdirTemp("", "dlp-scan-")
    if err != nil {
        return nil, err
    }
    defer os.RemoveAll(tempDir)

    // 2. Download files from GCS to local directory
    // (Your GCS download logic here)
    // ...

    // 3. Load DLP configuration
    config, err := dlp.LoadConfig("")
    if err != nil {
        return nil, err
    }

    // 4. Scan the directory
    scanner := dlp.NewScanner(config)
    return scanner.ScanDirectory(ctx, tempDir, bucketName, fileCount)
}
```

### Custom Configuration

```go
package main

import (
    "github.com/ca-risken/common/pkg/dlp"
)

func main() {
    // Load custom configuration file
    config, err := dlp.LoadConfig("/path/to/custom/dlp.yaml")
    if err != nil {
        log.Fatalf("Failed to load custom config: %v", err)
    }

    scanner := dlp.NewScanner(config)
    // ... use scanner
}
```

## Configuration

### DLP Configuration File (dlp.yaml)

```yaml
# Scanning limits
max_scan_files: 1000
max_scan_size_mb: 100
max_single_file_size_mb: 5
max_matches_per_finding: 5

# File patterns to exclude from scanning
exclude_file_patterns:
  - ".html"
  - ".css"
  - ".js"
  - "node_modules"

# Optional: Custom fingerprint file path
fingerprint_file_path: "yaml/fingerprint.yaml"

# Detection rules
rules:
  - name: "Email"
    description: "Email addresses"
    type: "PII"
    severity_thresholds:
      critical: 100
      high: 30
      medium: 5

  - name: "AWS Access Key ID"
    description: "AWS access key IDs"
    type: "Credential"
    severity_thresholds:
      critical: 1
```

## API Reference

### Types

#### Config
- `LoadConfig(configPath string) (*Config, error)`: Load configuration
- `GetMaxScanSizeBytes() int64`: Get maximum scan size in bytes
- `GetMaxSingleFileSizeBytes() int64`: Get maximum single file size in bytes
- `GetFingerprintFilePath() string`: Get fingerprint file path
- `CopyFingerprintFile(destDir string) (string, error)`: Copy fingerprint file
- `GetRule(patternName string) *Rule`: Get rule by pattern name

#### Rule
- `CalculateSeverity(matchCount int) string`: Calculate severity based on match count
- `IsApplicableToFile(fileName string, fileSizeBytes int64) bool`: Check if rule applies to file

#### Scanner
- `NewScanner(config *Config) *Scanner`: Create new scanner
- `ScanDirectory(ctx context.Context, dir string, resourceName string, totalFiles int) (*ScanResult, error)`: Scan directory

### Constants

```go
const (
    SeverityLow      = "LOW"
    SeverityMedium   = "MEDIUM"
    SeverityHigh     = "HIGH"
    SeverityCritical = "CRITICAL"
)
```

## Requirements

- Go 1.21 or later
- `hawk_scanner` command must be available in PATH

## License

Same as the parent RISKEN project.
