package log

import (
	"context"
	"io"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/memory"
)

var defaultLogger = NewLoggerWithParameters(os.Stdout, JSON, InfoLevel)

type Format string

const (
	CLI     Format = "cli"
	Discard Format = "discard"
	JSON    Format = "json"
	Memory  Format = "memory"
)

type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
)

func Debug(msg string, fields Fields) {
	defaultLogger.Debug(context.Background(), msg, fields)
}

func Entries() []*log.Entry {
	return defaultLogger.Entries()
}

func Error(err error, msg string, fields Fields) {
	defaultLogger.Error(context.Background(), err, msg, fields)
}

func Info(msg string, fields Fields) {
	defaultLogger.Info(context.Background(), msg, fields)
}

func NewLogger() Logger {
	return Logger{
		l: log.Logger{
			Level:   defaultLogger.l.Level,
			Handler: defaultLogger.l.Handler,
		},
	}
}

func NewLoggerWithParameters(w io.Writer, format Format, level Level) Logger {
	return Logger{
		l:      newApexLogger(w, format, level),
		w:      w,
		format: format,
		level:  level,
	}
}

func SetFormat(format Format) {
	defaultLogger.SetFormat(format)
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

func SetWriter(w io.Writer) {
	defaultLogger.SetWriter(w)
}

type Fields map[string]interface{}

func (fs Fields) Fields() log.Fields {
	return log.Fields(fs)
}

func (fs Fields) addContext(ctx context.Context) Fields {
	if fs == nil {
		fs = make(map[string]interface{})
	}
	// scope := cloud.GetScope(ctx)
	// if scope.RequestID != "" {
	// 	fs["request_id"] = scope.RequestID
	// }
	// if scope.CID != 0 {
	// 	fs["cid"] = scope.CID
	// }
	return fs
}

type Logger struct {
	l      log.Logger
	w      io.Writer
	format Format
	level  Level
}

func (l Logger) Debug(ctx context.Context, msg string, fields Fields) {
	fields = fields.addContext(ctx)
	l.l.WithFields(fields).Debug(msg)
}

func (l Logger) Entries() []*log.Entry {
	h := l.l.Handler.(*memory.Handler)
	return h.Entries
}

func (l Logger) Error(ctx context.Context, err error, msg string, fields Fields) {
	fields = fields.addContext(ctx)
	l.l.WithError(err).WithFields(fields).Error(msg)
}

func (l Logger) Info(ctx context.Context, msg string, fields Fields) {
	fields = fields.addContext(ctx)
	l.l.WithFields(fields).Info(msg)
}

func (l *Logger) SetFormat(format Format) {
	l.format = format
	l.l = newApexLogger(l.w, l.format, l.level)
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
	l.l = newApexLogger(l.w, l.format, l.level)
}

func (l *Logger) SetWriter(w io.Writer) {
	l.w = w
	l.l = newApexLogger(l.w, l.format, l.level)
}

func newApexLogger(w io.Writer, format Format, level Level) log.Logger {
	apexLogger := log.Logger{Level: log.InfoLevel}

	switch format {
	case CLI:
		apexLogger.Handler = cli.New(w)
	case Discard:
		apexLogger.Handler = discard.New()
	case Memory:
		apexLogger.Handler = memory.New()
	default:
		apexLogger.Handler = json.New(w)
	}

	if level == DebugLevel {
		apexLogger.Level = log.DebugLevel
	}

	return apexLogger
}
