package log

import (
	"fmt"
	"io"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
)

type Handler string

const (
	JSONHandler Handler = "json"
	CLIHandler  Handler = "cli"

	DefaultHandler = JSONHandler
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
)

type Entry struct {
	*log.Entry
}

type config struct {
	handler Handler
	level   log.Level
	w       io.Writer
}

type Option func(cfg *config)

func WithLogLevel(lvl LogLevel) Option {
	return func(cfg *config) {
		switch lvl {
		case DebugLevel:
			cfg.level = log.DebugLevel
		case InfoLevel:
			cfg.level = log.InfoLevel
		default:
			err := fmt.Errorf("invalid log level: %v", lvl)
			panic(err)
		}
	}
}

func WithHandler(h Handler) Option {
	return func(cfg *config) {
		switch h {
		case JSONHandler:
			cfg.handler = JSONHandler
		case CLIHandler:
			cfg.handler = CLIHandler
		default:
			err := fmt.Errorf("invalid handler: %v", h)
			panic(err)
		}
	}
}

func NewLogger(opts ...Option) Logger {
	cfg := config{
		handler: DefaultHandler,
		level:   log.InfoLevel,
		w:       os.Stderr,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	l := Logger{
		l: log.Logger{
			Level: cfg.level,
		},
	}
	switch cfg.handler {
	case CLIHandler:
		l.l.Handler = cli.New(cfg.w)
	default:
		l.l.Handler = json.New(cfg.w)
	}
	return l
}

type Logger struct {
	l  log.Logger
	id string
}

func (l Logger) Debug(msg string) {
	l.l.Debug(msg)
}

func (l Logger) Info(msg string) {
	l.l.Info(msg)
}

func (l Logger) WithRequestID(id string) Logger {
	l.id = id
	return l
}

func (l Logger) WithError(err error) *Entry {
	apexEntry := l.l.WithError(err)
	if l.id != "" {
		apexEntry = apexEntry.WithField("request_id", l.id)
	}
	return &Entry{Entry: apexEntry}
}

func (l Logger) WithField(key string, value interface{}) *Entry {
	apexEntry := l.l.WithField(key, value)
	if l.id != "" {
		apexEntry = apexEntry.WithField("request_id", l.id)
	}
	return &Entry{Entry: apexEntry}
}
