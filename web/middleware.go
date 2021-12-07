package web

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/log"
)

type EndpointMiddleware func(EndpointFunc) EndpointFunc

func LoggingMiddleware(next EndpointFunc) EndpointFunc {
	return func(ctx cloud.Context, request interface{}) (response interface{}, err *cloud.Error) {
		defer func(begin time.Time) {
			l := log.NewLogger()
			l.Info(ctx.Ctx, "completed request", log.Fields{
				"took":   time.Since(begin).Seconds(),
				"path":   ctx.Request.URL.Path,
				"method": ctx.Request.Method,
				"host":   ctx.Request.Host,
				"proto":  ctx.Request.Proto,
				"ua":     ctx.Request.UserAgent(),
			})
		}(time.Now())
		return next(ctx, request)
	}
}

// func MonitoringMiddleware(next EndpointFunc) EndpointFunc {
// 	return func(ctx cloud.Context, request interface{}) (response interface{}, err error) {
// 		defer func(begin time.Time) {
// 			sensor := instrumentation.NewSensor()
// 			sensor.Incr(ctx, "responses", nil)
// 			sensor.Timing(ctx, "latency", time.Since(begin), nil)
// 		}(time.Now())
// 		return next(ctx, request)
// 	}
// }

type ServerMiddleware func(next http.Handler) http.Handler

func HTTPOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			w.WriteHeader(http.StatusUpgradeRequired)
			w.Header().Set("Upgrade", "TLS/1.3, HTTP/1.1")
			w.Write([]byte("HTTP requests are not allowed for this service"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func H18(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			io.Copy(ioutil.Discard, r.Body)
			r.Body.Close()
		}()
		next.ServeHTTP(w, r)
	})
}
