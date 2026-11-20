package dlp

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

//go:embed yaml/dlp.yaml
var embeddedDLPYaml embed.FS

const defaultDLPConfigFile = "yaml/dlp.yaml"

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

// GetMaxSingleFileSizeBytes returns the maximum single file size in bytes
func (c *Config) GetMaxSingleFileSizeBytes() int64 {
	if c.MaxSingleFileSizeMB <= 0 {
		return 5 * 1024 * 1024 // Default 5MB if invalid
	}
	return int64(c.MaxSingleFileSizeMB) * 1024 * 1024
}
