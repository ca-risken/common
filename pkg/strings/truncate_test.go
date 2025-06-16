package strings

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		suffix    string
		want      string
	}{
		{
			name:      "no truncation needed",
			input:     "hello",
			maxLength: 10,
			suffix:    "...",
			want:      "hello",
		},
		{
			name:      "exact length match",
			input:     "hello",
			maxLength: 5,
			suffix:    "...",
			want:      "hello",
		},
		{
			name:      "truncation with ellipsis",
			input:     "hello world",
			maxLength: 5,
			suffix:    "...",
			want:      "hello...",
		},
		{
			name:      "truncation without suffix",
			input:     "hello world",
			maxLength: 5,
			suffix:    "",
			want:      "hello",
		},
		{
			name:      "empty string",
			input:     "",
			maxLength: 5,
			suffix:    "...",
			want:      "",
		},
		{
			name:      "japanese characters",
			input:     "こんにちは世界",
			maxLength: 5,
			suffix:    "...",
			want:      "こんにちは...",
		},
		{
			name:      "mixed ascii and japanese",
			input:     "hello世界",
			maxLength: 6,
			suffix:    "...",
			want:      "hello世...",
		},
		{
			name:      "zero max length",
			input:     "hello",
			maxLength: 0,
			suffix:    "...",
			want:      "...",
		},
		{
			name:      "custom suffix",
			input:     "long text here",
			maxLength: 4,
			suffix:    "…",
			want:      "long…",
		},
		{
			name:      "emoji characters",
			input:     "😀😂🎉🚀",
			maxLength: 2,
			suffix:    "...",
			want:      "😀😂...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateString(tt.input, tt.maxLength, tt.suffix)
			if got != tt.want {
				t.Errorf("TruncateString() = %v, want %v\n%s", got, tt.want, cmp.Diff(tt.want, got))
			}
		})
	}
}
