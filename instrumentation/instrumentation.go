package instrumentation

import (
	"context"
	"time"
)

var defaultSensor = Sensor{handler: discardHandler{}}

type Handler interface {
	HandleIncr(IncrEvent)
	HandleTiming(TimingEvent)
}

type IncrEvent struct {
	Key  string
	Tags Tags
}

type Sensor struct {
	handler Handler
}

func NewSensor() Sensor {
	return Sensor{
		handler: defaultSensor.handler,
	}
}

func (s Sensor) Incr(ctx context.Context, key string, tags Tags) {
	tags = tags.addContext(ctx)
	event := IncrEvent{
		Key:  key,
		Tags: tags,
	}
	s.handler.HandleIncr(event)
}

func (s Sensor) Timing(ctx context.Context, key string, duration time.Duration, tags Tags) {
	tags = tags.addContext(ctx)
	event := TimingEvent{
		Duration: duration,
		Key:      key,
		Tags:     tags,
	}
	s.handler.HandleTiming(event)
}

type Tags map[string]string

func (ts Tags) addContext(ctx context.Context) Tags {
	// if ts == nil {
	// 	ts = make(map[string]string)
	// }
	// scope := cloud.GetScope(ctx)
	// if scope.Request != nil {
	// 	ts["path"] = scope.Request.URL.Path
	// 	ts["ua"] = scope.Request.UserAgent()
	// 	ts["method"] = scope.Request.Method
	// }
	return ts
}

type TimingEvent struct {
	Duration time.Duration
	Key      string
	Tags     Tags
}

func Incr(key string, tags Tags) {
	defaultSensor.Incr(context.Background(), key, tags)
}

func SetHandler(handler Handler) {
	defaultSensor.handler = handler
}

func Timing(key string, duration time.Duration, tags Tags) {
	defaultSensor.Timing(context.Background(), key, duration, tags)
}

type discardHandler struct{}

func (h discardHandler) HandleIncr(_ IncrEvent) {
}

func (h discardHandler) HandleTiming(_ TimingEvent) {}
