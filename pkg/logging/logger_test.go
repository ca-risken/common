package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	logger := NewLogger()
	cases := []struct {
		name       string
		input      string
		wantPrefix string
		loops      int
	}{
		{
			name:       "OK-100",
			input:      "prefix",
			wantPrefix: "prefix-",
			loops:      100,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotIDs := []string{}
			for i := 0; i < c.loops; i++ {
				got, err := logger.GenerateRequestID(c.input)
				if err != nil {
					t.Fatalf("Unexpected GenerateRequestID error: err=%+v", err)
				}
				if !strings.HasPrefix(got, c.wantPrefix) {
					t.Fatalf("Unexpected not match ID format: want=%s, got=%s", c.wantPrefix, got)
				}
				// fmt.Println(got)
				gotIDs = append(gotIDs, got)
			}
			for _, target := range gotIDs {
				sameCnt := 0
				for _, id := range gotIDs {
					if target == id {
						sameCnt++
					}
					if sameCnt > 1 {
						t.Fatalf("Detect same id pattern on the %d loops, ID=%s", c.loops, target)
					}
				}
			}
		})
	}
}

func TestMustNotify(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", notifyKey, "true"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.SetOutput(buf)
			logger.MustNotify(InfoLevel, c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(inclued keyword)=%s, got=%s", key, logged)
				}

			}
		})
	}
}

func TestMustNotifyf(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", notifyKey, "true"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.SetOutput(buf)
			logger.MustNotifyf(InfoLevel, "%s", c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(inclued keyword)=%s, got=%s", key, logged)
				}

			}
		})
	}
}
