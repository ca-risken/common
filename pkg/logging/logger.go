package logging

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Level type
type Level uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel

	notifyKey       = "notify"
	FieldKeyTraceID = "dd.trace_id"
	FieldKeySpanID  = "dd.span_id"
)

// Logger interface
type Logger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})

	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Panic(ctx context.Context, args ...interface{})

	Level(level Level)
	Output(w io.Writer)
	GenerateRequestID(prefix string) (string, error)
	Notify(ctx context.Context, level Level, args ...interface{})
	Notifyf(ctx context.Context, level Level, format string, args ...interface{})
	WithItems(ctx context.Context, level Level, fields map[string]interface{}, args ...interface{})
	WithItemsf(ctx context.Context, level Level, fields map[string]interface{}, format string, args ...interface{})
}

type AppLogger struct {
	*logrus.Logger
}

func NewLogger() Logger {
	l := logrus.New()
	l.SetOutput(os.Stderr)
	l.SetFormatter(&logrus.JSONFormatter{})
	l.Level = logrus.InfoLevel
	return &AppLogger{l}
}

func (a *AppLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	a.logWithTracef(ctx, DebugLevel, nil, format, args...)
}

func (a *AppLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	a.logWithTracef(ctx, InfoLevel, nil, format, args...)
}

func (a *AppLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	a.logWithTracef(ctx, WarnLevel, nil, format, args...)
}

func (a *AppLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	a.logWithTracef(ctx, ErrorLevel, nil, format, args...)
}

func (a *AppLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	a.logWithTracef(ctx, FatalLevel, nil, format, args...)
	a.Exit(1)
}

func (a *AppLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	a.logWithTracef(ctx, PanicLevel, nil, format, args...)
}

func (a *AppLogger) Debug(ctx context.Context, args ...interface{}) {
	a.logWithTrace(ctx, DebugLevel, nil, args...)
}

func (a *AppLogger) Info(ctx context.Context, args ...interface{}) {
	a.logWithTrace(ctx, InfoLevel, nil, args...)
}

func (a *AppLogger) Warn(ctx context.Context, args ...interface{}) {
	a.logWithTrace(ctx, WarnLevel, nil, args...)
}

func (a *AppLogger) Error(ctx context.Context, args ...interface{}) {
	a.logWithTrace(ctx, ErrorLevel, nil, args...)
}

func (a *AppLogger) Fatal(ctx context.Context, args ...interface{}) {
	a.logWithTrace(ctx, FatalLevel, nil, args...)
	a.Exit(1)
}

func (a *AppLogger) Panic(ctx context.Context, args ...interface{}) {
	a.logWithTrace(ctx, PanicLevel, nil, args...)
}

func parseLogrusLevel(l Level) logrus.Level {
	switch l {
	case PanicLevel:
		return logrus.PanicLevel
	case FatalLevel:
		return logrus.FatalLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case WarnLevel:
		return logrus.WarnLevel
	case InfoLevel:
		return logrus.InfoLevel
	case DebugLevel:
		return logrus.DebugLevel
	case TraceLevel:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

func (a *AppLogger) Level(level Level) {
	a.SetLevel(parseLogrusLevel(level))
}

func (a *AppLogger) Output(w io.Writer) {
	a.SetOutput(w)
}

func (a *AppLogger) GenerateRequestID(prefix string) (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", prefix, u.String()), nil
}

func (a *AppLogger) Notify(ctx context.Context, level Level, args ...interface{}) {
	m := map[string]interface{}{notifyKey: true}
	a.logWithTrace(ctx, level, m, args...)
}

func (a *AppLogger) Notifyf(ctx context.Context, level Level, format string, args ...interface{}) {
	m := map[string]interface{}{notifyKey: true}
	a.logWithTracef(ctx, level, m, format, args...)
}

func (a *AppLogger) WithItems(ctx context.Context, level Level, fields map[string]interface{}, args ...interface{}) {
	a.logWithTrace(ctx, level, fields, args...)
}

func (a *AppLogger) WithItemsf(ctx context.Context, level Level, fields map[string]interface{}, format string, args ...interface{}) {
	a.logWithTracef(ctx, level, fields, format, args...)
}

func (a *AppLogger) logWithTracef(ctx context.Context, level Level, fields map[string]interface{}, format string, args ...interface{}) {
	f := logrus.Fields{}
	for k, v := range fields {
		f[k] = v
	}
	span, ok := tracer.SpanFromContext(ctx)
	if ok {
		f[FieldKeyTraceID] = span.Context().TraceID()
		f[FieldKeySpanID] = span.Context().SpanID()
	}
	a.WithFields(f).Logf(parseLogrusLevel(level), format, args...)
}

func (a *AppLogger) logWithTrace(ctx context.Context, level Level, fields map[string]interface{}, args ...interface{}) {
	f := logrus.Fields{}
	for k, v := range fields {
		f[k] = v
	}
	span, ok := tracer.SpanFromContext(ctx)
	if ok {
		f[FieldKeyTraceID] = span.Context().TraceID()
		f[FieldKeySpanID] = span.Context().SpanID()
	}
	a.WithFields(f).Log(parseLogrusLevel(level), args...)
}
