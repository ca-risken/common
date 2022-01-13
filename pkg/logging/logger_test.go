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

func TestDebug(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "debug"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Debug(c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestDebugf(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "debug"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Debugf("%s", c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestInfo(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "info"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Info(c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestInfof(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "info"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Infof("%s", c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestWarn(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "warning"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Warn(c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestWarnf(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "warning"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Warnf("%s", c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestWarning(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "warning"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Warning(c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestWarningf(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "warning"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Warningf("%s", c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestError(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "error"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Error(c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "OK",
			input: "something",
			want:  []string{"something", "error"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Errorf("%s", c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestNotify(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		inputLevel Level
		want       []string
	}{
		{
			name:       "OK",
			input:      "something",
			inputLevel: InfoLevel,
			want:       []string{"something", "info", notifyKey, "true"},
		},
		{
			name:       "OK debug",
			input:      "something",
			inputLevel: DebugLevel,
			want:       []string{"something", "debug", notifyKey, "true"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Notify(c.inputLevel, c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestNotifyf(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		inputLevel Level
		want       []string
	}{
		{
			name:       "OK",
			input:      "something",
			inputLevel: InfoLevel,
			want:       []string{"something", "info", notifyKey, "true"},
		},
		{
			name:       "OK debug",
			input:      "something",
			inputLevel: DebugLevel,
			want:       []string{"something", "debug", notifyKey, "true"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.Notifyf(c.inputLevel, "%s", c.input)
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}
func TestWithItems(t *testing.T) {
	cases := []struct {
		name       string
		input      map[string]interface{}
		inputLevel Level
		want       []string
	}{
		{
			name: "OK single item",
			input: map[string]interface{}{
				"key1": "value1",
			},
			inputLevel: InfoLevel,
			want:       []string{"info", "key1", "value1"},
		},
		{
			name: "OK multiple items",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			inputLevel: DebugLevel,
			want:       []string{"debug", "key1", "value1", "key2", "value2"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.WithItems(c.inputLevel, c.input, "test")
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}

func TestWithItemsf(t *testing.T) {
	cases := []struct {
		name       string
		input      map[string]interface{}
		inputLevel Level
		want       []string
	}{
		{
			name: "OK single item",
			input: map[string]interface{}{
				"key1": "value1",
			},
			inputLevel: InfoLevel,
			want:       []string{"info", "key1", "value1"},
		},
		{
			name: "OK multiple items",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			inputLevel: DebugLevel,
			want:       []string{"debug", "key1", "value1", "key2", "value2"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			logger.WithItemsf(c.inputLevel, c.input, "%s", "test")
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}
