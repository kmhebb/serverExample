package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/log"
)

const (
	DefaultCircuitBreakerThreshold int           = 1
	DefaultCircuitBreakerWait      time.Duration = 20 * time.Second
)

type CircuitBreaker func() error

type ServerState string

const (
	LiveState     ServerState = "live"
	ReadyState    ServerState = "ready"
	StartingState ServerState = "starting"
	StoppingState ServerState = "stopping"
)

func NewUnstartedServer(addr string) *Server {
	return newServer(addr)
}

func NewServer(addr string) *Server {
	srv := newServer(addr)
	srv.Start()
	return srv
}

func newServer(addr string) *Server {
	srv := &Server{
		Addr:                    addr,
		CircuitBreakerWait:      DefaultCircuitBreakerWait,
		circuitBroken:           0,
		circuitBreakerThreshold: DefaultCircuitBreakerThreshold,
		router:                  mux.NewRouter(),
		started:                 time.Now(),
	}

	srv.srv = &http.Server{
		Addr:    addr,
		Handler: srv,
	}

	srv.api = srv.router.PathPrefix("/api").Subrouter()
	srv.api.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := cloud.NewContext(r)
		logger := log.NewLogger()
		logger.Info(ctx.Ctx, "route not found", log.Fields{"path": r.URL.Path})
	})

	return srv
}

type statusData struct {
	State  ServerState `json:"state"`
	Uptime string      `json:"uptime"`
}

type Server struct {
	Addr               string
	CircuitBreakerWait time.Duration

	api                     *mux.Router
	circuitBreakers         []CircuitBreaker
	circuitBroken           int
	circuitBreakerThreshold int
	err                     error
	router                  *mux.Router
	srv                     *http.Server
	started                 time.Time
	state                   ServerState
}

func (srv *Server) AddCircuitBreakers(cbs ...CircuitBreaker) {
	srv.circuitBreakers = append(srv.circuitBreakers, cbs...)
}

func (srv *Server) Handle(path string, h http.Handler) {
	srv.api.Handle(path, h)
}

func (srv *Server) Listen() {
	if srv.state != StartingState && srv.state != LiveState {
		return
	}
	srv.state = ReadyState
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health/GetStatus" {
		w.Header().Set("Content-Type", ContentTypeJSON)
		json.NewEncoder(w).Encode(statusData{
			State:  srv.state,
			Uptime: time.Since(srv.started).Truncate(time.Millisecond).String(),
		})
		return
	}

	if srv.state != ReadyState {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	srv.router.ServeHTTP(w, r)
}

func (srv *Server) Start() {
	if srv.state != "" {
		return
	}
	srv.state = StartingState

	go srv.circuitBreak()

	go func() {
		srv.err = srv.srv.ListenAndServe()
	}()
}

func (srv *Server) Stop() error {
	if srv.state != ReadyState && srv.state != LiveState {
		return nil
	}
	srv.state = StoppingState
	if err := srv.srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("web/Server.Stop: %w", err)
	}
	return nil
}

func (srv *Server) Use(mws ...ServerMiddleware) {
	for _, mw := range mws {
		srv.api.Use(mux.MiddlewareFunc(mw))
	}
}

func (srv *Server) circuitBreak() {
	circuitBroken := 0
	for {
		<-time.After(srv.CircuitBreakerWait)

		inError := srv.state == LiveState
		switch {
		case inError:

			// If we're already in error, we need an all clear from every circuit
			// breaker before we restore a ready status.
			allClear := true
			for _, f := range srv.circuitBreakers {
				err := f()
				allClear = allClear && (err == nil)
			}

			// If we don't get the all clear, we increment our counter tracking how many
			// times we've tried to restore functionality.
			if !allClear {
				circuitBroken++
				log.Info("circuits in error", log.Fields{
					"attempt":   circuitBroken,
					"threshold": srv.circuitBreakerThreshold,
				})
				// If we've hit our limit, we shutdown the server
				if circuitBroken == srv.circuitBreakerThreshold {
					log.Info("circuit breaker threshold reached. stopping server.", nil)
					// TODO: Send an error to sentry
					srv.Stop()
					return
				}
				continue
			}

			// If everything looks like it's back to normal, we start listening and
			// reset our counter
			log.Info("circuit breakers reset.", nil)
			srv.state = ReadyState
			circuitBroken = 0
		default:
			// If we're not in error, we just run our circuit breakers looking for the
			// first error.
			for _, f := range srv.circuitBreakers {
				if err := f(); err != nil {
					srv.state = LiveState
				}
			}
		}
	}
}
