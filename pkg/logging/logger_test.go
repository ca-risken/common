package logging

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
		name    string
		input   string
		traceOn bool
		want    []string
	}{
		{
			name:    "OK",
			input:   "something",
			traceOn: false,
			want:    []string{"something", "debug"},
		},
		{
			name:    "OK trace",
			input:   "something",
			traceOn: true,
			want:    []string{"something", "debug", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Debug(ctx, c.input)
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
		name    string
		input   string
		traceOn bool
		want    []string
	}{
		{
			name:    "OK",
			input:   "something",
			traceOn: false,
			want:    []string{"something", "debug"},
		},
		{
			name:    "OK trace",
			input:   "something",
			traceOn: true,
			want:    []string{"something", "debug", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Debugf(ctx, "%s", c.input)
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
		name    string
		input   string
		traceOn bool
		want    []string
	}{
		{
			name:    "OK",
			input:   "something",
			traceOn: false,
			want:    []string{"something", "info"},
		},
		{
			name:    "OK trace",
			input:   "something",
			traceOn: true,
			want:    []string{"something", "info", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Info(ctx, c.input)
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
		name    string
		input   string
		traceOn bool
		want    []string
	}{
		{
			name:    "OK",
			input:   "something",
			traceOn: false,
			want:    []string{"something", "info"},
		},
		{
			name:    "OK trace",
			input:   "something",
			traceOn: true,
			want:    []string{"something", "info", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Infof(ctx, "%s", c.input)
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
		name    string
		input   string
		traceOn bool
		want    []string
	}{
		{
			name:    "OK",
			input:   "something",
			traceOn: false,
			want:    []string{"something", "warning"},
		},
		{
			name:    "OK trace",
			input:   "something",
			traceOn: true,
			want:    []string{"something", "warning", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Warn(ctx, c.input)
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
		name    string
		input   string
		traceOn bool
		want    []string
	}{
		{
			name:    "OK",
			input:   "something",
			traceOn: false,
			want:    []string{"something", "warning"},
		},
		{
			name:    "OK trace",
			input:   "something",
			traceOn: true,
			want:    []string{"something", "warning", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Warnf(ctx, "%s", c.input)
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
		name    string
		input   string
		traceOn bool
		want    []string
	}{
		{
			name:    "OK",
			input:   "something",
			traceOn: false,
			want:    []string{"something", "error"},
		},
		{
			name:    "OK trace",
			input:   "something",
			traceOn: true,
			want:    []string{"something", "error", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Error(ctx, c.input)
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
		name    string
		input   string
		traceOn bool
		want    []string
	}{
		{
			name:    "OK",
			input:   "something",
			traceOn: false,
			want:    []string{"something", "error"},
		},
		{
			name:    "OK trace",
			input:   "something",
			traceOn: true,
			want:    []string{"something", "error", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Errorf(ctx, "%s", c.input)
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
		traceOn    bool
		want       []string
	}{
		{
			name:       "OK",
			input:      "something",
			inputLevel: InfoLevel,
			traceOn:    false,
			want:       []string{"something", "info", notifyKey, "true"},
		},
		{
			name:       "OK debug",
			input:      "something",
			inputLevel: DebugLevel,
			traceOn:    false,
			want:       []string{"something", "debug", notifyKey, "true"},
		},
		{
			name:       "OK trace",
			input:      "something",
			inputLevel: InfoLevel,
			traceOn:    true,
			want:       []string{"something", "info", notifyKey, "true", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Notify(ctx, c.inputLevel, c.input)
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
		traceOn    bool
		want       []string
	}{
		{
			name:       "OK",
			input:      "something",
			inputLevel: InfoLevel,
			traceOn:    false,
			want:       []string{"something", "info", notifyKey, "true"},
		},
		{
			name:       "OK debug",
			input:      "something",
			inputLevel: DebugLevel,
			traceOn:    false,
			want:       []string{"something", "debug", notifyKey, "true"},
		},
		{
			name:       "OK trace",
			input:      "something",
			inputLevel: InfoLevel,
			traceOn:    true,
			want:       []string{"something", "info", notifyKey, "true", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.Notifyf(ctx, c.inputLevel, "%s", c.input)
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
		traceOn    bool
		want       []string
	}{
		{
			name: "OK single item",
			input: map[string]interface{}{
				"key1": "value1",
			},
			inputLevel: InfoLevel,
			traceOn:    false,
			want:       []string{"info", "key1", "value1"},
		},
		{
			name: "OK multiple items",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			inputLevel: DebugLevel,
			traceOn:    false,
			want:       []string{"debug", "key1", "value1", "key2", "value2"},
		},
		{
			name: "OK trace",
			input: map[string]interface{}{
				"key1": "value1",
			},
			inputLevel: InfoLevel,
			traceOn:    true,
			want:       []string{"info", "key1", "value1", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.WithItems(ctx, c.inputLevel, c.input, "test")
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
		traceOn    bool
		want       []string
	}{
		{
			name: "OK single item",
			input: map[string]interface{}{
				"key1": "value1",
			},
			inputLevel: InfoLevel,
			traceOn:    false,
			want:       []string{"info", "key1", "value1"},
		},
		{
			name: "OK multiple items",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			inputLevel: DebugLevel,
			traceOn:    false,
			want:       []string{"debug", "key1", "value1", "key2", "value2"},
		},
		{
			name: "OK trace",
			input: map[string]interface{}{
				"key1": "value1",
			},
			inputLevel: InfoLevel,
			traceOn:    true,
			want:       []string{"info", "key1", "value1", "dd.trace_id", "dd.span_id"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewLogger()
			logger.Output(buf)
			logger.Level(DebugLevel)
			ctx := context.Background()
			if c.traceOn {
				span := tracer.StartSpan("test span")
				defer span.Finish()
				ctx = tracer.ContextWithSpan(ctx, span)
			}
			logger.WithItemsf(ctx, c.inputLevel, c.input, "%s", "test")
			logged := buf.String()
			for _, key := range c.want {
				if !strings.Contains(logged, key) {
					t.Fatalf("Unexpected log: want(keyword)=%s, got=%s", key, logged)
				}
			}
		})
	}
}
