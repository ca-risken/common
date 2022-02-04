package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
)

// Logger interface
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Level(level Level)
	Output(w io.Writer)
	GenerateRequestID(prefix string) (string, error)
	Notify(level Level, args ...interface{})
	Notifyf(level Level, format string, args ...interface{})
	WithItems(level Level, fields map[string]interface{}, args ...interface{})
	WithItemsf(level Level, fields map[string]interface{}, format string, args ...interface{})
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

func (a *AppLogger) Debugf(format string, args ...interface{}) {
	a.Logf(logrus.DebugLevel, format, args...)
}

func (a *AppLogger) Infof(format string, args ...interface{}) {
	a.Logf(logrus.InfoLevel, format, args...)
}

func (a *AppLogger) Warnf(format string, args ...interface{}) {
	a.Logf(logrus.WarnLevel, format, args...)
}

func (a *AppLogger) Warningf(format string, args ...interface{}) {
	a.Warnf(format, args...)
}

func (a *AppLogger) Errorf(format string, args ...interface{}) {
	a.Logf(logrus.ErrorLevel, format, args...)
}

func (a *AppLogger) Fatalf(format string, args ...interface{}) {
	a.Logf(logrus.FatalLevel, format, args...)
	a.Exit(1)
}

func (a *AppLogger) Panicf(format string, args ...interface{}) {
	a.Logf(logrus.PanicLevel, format, args...)
}

func (a *AppLogger) Debug(args ...interface{}) {
	a.Log(logrus.DebugLevel, args...)
}

func (a *AppLogger) Info(args ...interface{}) {
	a.Log(logrus.InfoLevel, args...)
}

func (a *AppLogger) Warn(args ...interface{}) {
	a.Log(logrus.WarnLevel, args...)
}

func (a *AppLogger) Warning(args ...interface{}) {
	a.Warn(args...)
}

func (a *AppLogger) Error(args ...interface{}) {
	a.Log(logrus.ErrorLevel, args...)
}

func (a *AppLogger) Fatal(args ...interface{}) {
	a.Log(logrus.FatalLevel, args...)
	a.Exit(1)
}

func (a *AppLogger) Panic(args ...interface{}) {
	a.Log(logrus.PanicLevel, args...)
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

const (
	notifyKey = "notify"
)

func (a *AppLogger) Notify(level Level, args ...interface{}) {
	a.WithFields(logrus.Fields{
		"notify": true,
	}).Log(parseLogrusLevel(level), args...)
}

func (a *AppLogger) Notifyf(level Level, format string, args ...interface{}) {
	a.WithFields(logrus.Fields{
		"notify": true,
	}).Logf(parseLogrusLevel(level), format, args...)
}

func (a *AppLogger) WithItems(level Level, fields map[string]interface{}, args ...interface{}) {
	f := logrus.Fields{}
	for k, v := range fields {
		f[k] = v
	}
	a.WithFields(f).Log(parseLogrusLevel(level), args...)
}

func (a *AppLogger) WithItemsf(level Level, fields map[string]interface{}, format string, args ...interface{}) {
	f := logrus.Fields{}
	for k, v := range fields {
		f[k] = v
	}
	a.WithFields(f).Logf(parseLogrusLevel(level), format, args...)
}