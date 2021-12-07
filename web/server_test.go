package web_test

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/kmhebb/serverExample/web"
)

type failNCircuitBreaker struct {
	n     int
	count int
}

func (cb *failNCircuitBreaker) run() error {
	if cb.count < cb.n {
		cb.count++
		return fmt.Errorf("fail")
	}
	return nil
}

func TestServer(t *testing.T) {
	srv := web.NewServer(":8080")
	srv.CircuitBreakerWait = time.Second

	srv.Start()
	time.Sleep(time.Second)

	curl(t, 200, "http://localhost:8080/health/GetStatus")
	curl(t, 503, "http://localhost:8080/api/test/Greet")

	srv.Handle("/test/Greet", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))

	curl(t, 200, "http://localhost:8080/health/GetStatus")
	curl(t, 503, "http://localhost:8080/api/test/Greet")

	srv.Listen()

	curl(t, 200, "http://localhost:8080/health/GetStatus")
	curl(t, 200, "http://localhost:8080/api/test/Greet")

	failTwice := &failNCircuitBreaker{n: 2}
	srv.AddCircuitBreakers(failTwice.run)

	<-time.After(time.Second)

	curl(t, 200, "http://localhost:8080/health/GetStatus")
	curl(t, 503, "http://localhost:8080/api/test/Greet")

	<-time.After(5 * time.Second)

	curl(t, 200, "http://localhost:8080/health/GetStatus")
	curl(t, 200, "http://localhost:8080/api/test/Greet")

	srv.Stop()
}

func curl(t *testing.T, code int, url string) {
	resp, err := http.Get(url)
	if err != nil {
		t.Log("unexpected error:", err)
		return
	}

	b, _ := httputil.DumpResponse(resp, true)
	t.Log(string(b))

	if resp.StatusCode != code {
		t.Errorf("expected code %d, got %d", code, resp.StatusCode)
	}
}
