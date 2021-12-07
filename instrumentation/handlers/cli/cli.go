package cli

import (
	"io"

	"github.com/kmhebb/serverExample/instrumentation"
	"github.com/kmhebb/serverExample/log"
)

func New(w io.Writer) Handler {
	return Handler{w: w}
}

type Handler struct {
	w io.Writer
}

func (h Handler) HandleIncr(evt instrumentation.IncrEvent) {
	fields := log.Fields{}
	for k, v := range evt.Tags {
		fields[k] = v
	}
	log.Info("incr", fields)
}

func (h Handler) HandleTiming(evt instrumentation.TimingEvent) {
	fields := log.Fields{"dur": evt.Duration}
	for k, v := range evt.Tags {
		fields[k] = v
	}
	log.Info("timing", fields)
}
