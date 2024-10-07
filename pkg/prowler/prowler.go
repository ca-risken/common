package prowler

import (
	"errors"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type ProwlerSetting struct {
	IgnorePlugin          []string                 `yaml:"ignorePlugin"`
	SpecificPluginSetting map[string]PluginSetting `yaml:"specificPluginSetting,omitempty" validate:"dive"`
}

type PluginSetting struct {
	Score                   *float32         `yaml:"score,omitempty"`
	SkipResourceNamePattern []string         `yaml:"skipResourceNamePattern,omitempty"`
	IgnoreMessagePattern    []string         `yaml:"ignoreMessagePattern,omitempty" validate:"dive,regexp"`
	Tags                    []string         `yaml:"tags,omitempty"`
	Recommend               *PluginRecommend `yaml:"recommend,omitempty"`
}

type PluginRecommend struct {
	Risk           *string `yaml:"risk,omitempty"`
	Recommendation *string `yaml:"recommendation,omitempty"`
}

func LoadProwlerSetting(path string) (*ProwlerSetting, error) {
	data, err := readProwlerSetting(path)
	if err != nil {
		return nil, err
	}

	setting, err := ParseProwlerSettingYaml(data)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func readProwlerSetting(path string) ([]byte, error) {
	var data []byte
	var err error
	if path == "" {
		return nil, errors.New("path is empty")
	}

	data, err = os.ReadFile(path) // Read from path
	if err != nil {
		return nil, err
	}
	return data, nil
}

func validateRegexp(fl validator.FieldLevel) bool {
	pattern, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	_, err := regexp.Compile(pattern)
	return err == nil
}

func ParseProwlerSettingYaml(data []byte) (*ProwlerSetting, error) {
	var setting ProwlerSetting
	if err := yaml.Unmarshal(data, &setting); err != nil {
		return nil, err
	}

	// validate
	validate := validator.New()
	if err := validate.RegisterValidation("regexp", validateRegexp); err != nil {
		return nil, err
	}
	if err := validate.Struct(setting); err != nil {
		return nil, err
	}
	return &setting, nil
}

func (c *ProwlerSetting) IsIgnorePlugin(plugin string) bool {
	return slices.Contains(c.IgnorePlugin, plugin)
}

func (c *ProwlerSetting) IsSkipResourceNamePattern(plugin, resourceName, aliasResourceName string) bool {
	if c.SpecificPluginSetting[plugin].SkipResourceNamePattern == nil {
		return false
	}
	for _, pattern := range c.SpecificPluginSetting[plugin].SkipResourceNamePattern {
		if strings.Contains(resourceName, pattern) {
			return true
		}
		if aliasResourceName != "" && strings.Contains(aliasResourceName, pattern) {
			return true
		}
	}
	return false
}

func (c *ProwlerSetting) IsIgnoreMessagePattern(plugin string, messages []string) bool {
	if c.SpecificPluginSetting[plugin].IgnoreMessagePattern == nil {
		return false
	}
	for _, pattern := range c.SpecificPluginSetting[plugin].IgnoreMessagePattern {
		for _, message := range messages {
			matched, err := regexp.MatchString(pattern, message)
			if err != nil {
				return false
			}
			if matched {
				return true // match == ignore finding
			}
		}
	}
	return false
}

// Helper function: return pointer of a value
func Ptr[T any](v T) *T {
	return &v
}
