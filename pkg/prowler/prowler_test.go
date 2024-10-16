package prowler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	TEST_YAML = `
ignorePlugin:
  - plugin1
  - plugin2
specificPluginSetting:
  plugin3:
    score: 0.75
    skipResourceNamePattern:
      - "_test"
    ignoreMessagePattern:
      - 'Domain: .+ expires in (?:2[5-9]|[3-9]\d|\d{3,}) days'
    tags:
      - tag1
      - tag2
    recommend:
      risk: "High risk"
      recommendation: "Fix it"
`
)

func TestParseDefaultProwlerSetting(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *ProwlerSetting
		wantErr bool
	}{
		{
			name:  "Valid YAML",
			input: []byte(TEST_YAML),
			want: &ProwlerSetting{
				IgnorePlugin: []string{"plugin1", "plugin2"},
				SpecificPluginSetting: map[string]PluginSetting{
					"plugin3": {
						Score:                   Ptr(float32(0.75)),
						SkipResourceNamePattern: []string{"_test"},
						IgnoreMessagePattern:    []string{`Domain: .+ expires in (?:2[5-9]|[3-9]\d|\d{3,}) days`},
						Tags:                    []string{"tag1", "tag2"},
						Recommend: &PluginRecommend{
							Risk:           Ptr("High risk"),
							Recommendation: Ptr("Fix it"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid structure",
			input:   []byte("invalid: yaml: content:"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid YAML",
			input:   []byte("foo: [bar: baz}"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseProwlerSettingYaml(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseProwlerSettingYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, result); diff != "" {
					t.Errorf("parseProwlerSettingYaml() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestIsIgnorePlugin(t *testing.T) {
	type args struct {
		setting *ProwlerSetting
		plugin  string
	}
	tests := []struct {
		name  string
		input args
		want  bool
	}{
		{
			name: "Ignore plugin",
			input: args{
				setting: &ProwlerSetting{
					IgnorePlugin: []string{"plugin1", "plugin2"},
				},
				plugin: "plugin1",
			},
			want: true,
		},
		{
			name: "Not ignore plugin",
			input: args{
				setting: &ProwlerSetting{
					IgnorePlugin: []string{"plugin1", "plugin2"},
				},
				plugin: "plugin3",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.setting.IsIgnorePlugin(tt.input.plugin); got != tt.want {
				t.Errorf("IsIgnorePlugin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSkipResourceNamePattern(t *testing.T) {
	type args struct {
		setting           *ProwlerSetting
		plugin            string
		resourceName      string
		aliasResourceName string
	}
	tests := []struct {
		name  string
		input args
		want  bool
	}{
		{
			name: "Skip resource name pattern matches",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {
							SkipResourceNamePattern: []string{"ignore", "test"},
						},
					},
				},
				plugin:            "plugin1",
				resourceName:      "testResourceName",
				aliasResourceName: "alias",
			},
			want: true,
		},
		{
			name: "Skip alias name pattern matches",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {
							SkipResourceNamePattern: []string{"ignore", "test"},
						},
					},
				},
				plugin:            "plugin1",
				resourceName:      "ResourceName",
				aliasResourceName: "ignoreAlias",
			},
			want: true,
		},
		{
			name: "Skip resource name pattern does not match",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {
							SkipResourceNamePattern: []string{"ignore", "test"},
						},
					},
				},
				plugin:            "plugin1",
				resourceName:      "resourceName",
				aliasResourceName: "alias",
			},
			want: false,
		},
		{
			name: "No skip resource name pattern",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {},
					},
				},
				plugin:            "plugin1",
				resourceName:      "resourceName",
				aliasResourceName: "alias",
			},
			want: false,
		},
		{
			name: "Plugin not found",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{},
				},
				plugin:            "plugin1",
				resourceName:      "resourceName",
				aliasResourceName: "alias",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.setting.IsSkipResourceNamePattern(tt.input.plugin, tt.input.resourceName, tt.input.aliasResourceName); got != tt.want {
				t.Errorf("IsSkipResourceNamePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIgnoreMessagePattern(t *testing.T) {
	type args struct {
		setting  *ProwlerSetting
		plugin   string
		messages []string
	}
	tests := []struct {
		name  string
		input args
		want  bool
	}{
		{
			name: "Ignore message pattern matches(25 days)",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {
							IgnoreMessagePattern: []string{`Domain: .+ expires in (?:2[5-9]|[3-9]\d|\d{3,}) days`},
						},
					},
				},
				plugin:   "plugin1",
				messages: []string{"Domain: example.com expires in 25 days"},
			},
			want: true,
		},
		{
			name: "Ignore message pattern matches(30 days)",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {
							IgnoreMessagePattern: []string{`Domain: .+ expires in (?:2[5-9]|[3-9]\d|\d{3,}) days`},
						},
					},
				},
				plugin:   "plugin1",
				messages: []string{"Domain: example.com expires in 30 days"},
			},
			want: true,
		},
		{
			name: "Ignore message pattern matches(100 days)",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {
							IgnoreMessagePattern: []string{`Domain: .+ expires in (?:2[5-9]|[3-9]\d|\d{3,}) days`},
						},
					},
				},
				plugin:   "plugin1",
				messages: []string{"Domain: example.com expires in 100 days"},
			},
			want: true,
		},
		{
			name: "Ignore message pattern does not match",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {
							IgnoreMessagePattern: []string{`Domain: .+ expires in (?:2[5-9]|[3-9]\d|\d{3,}) days`},
						},
					},
				},
				plugin:   "plugin1",
				messages: []string{"Domain: example.com expires in 24 days"},
			},
			want: false,
		},
		{
			name: "No ignore message pattern",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{
						"plugin1": {},
					},
				},
				plugin:   "plugin1",
				messages: []string{"Domain: example.com expires in 300 days"},
			},
			want: false,
		},
		{
			name: "Plugin not found",
			input: args{
				setting: &ProwlerSetting{
					SpecificPluginSetting: map[string]PluginSetting{},
				},
				plugin:   "plugin1",
				messages: []string{"Domain: example.com expires in 300 days"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.setting.IsIgnoreMessagePattern(tt.input.plugin, tt.input.messages); got != tt.want {
				t.Errorf("IsIgnoreMessagePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
