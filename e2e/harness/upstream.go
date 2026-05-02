//go:build e2e

package harness

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"
)

// Echo is the JSON shape returned by the mock upstream for every received
// request. Subtests deserialize this to assert the gateway forwarded
// correctly (preserved method, path, body, injected headers like
// X-Request-Id and Host).
type Echo struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Host    string              `json:"host"` // value the upstream saw on Request.Host
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

// Upstream is the in-process mock service altair forwards to. Special
// behaviors are triggered by URL path substrings so subtests can stay
// stateless:
//
//	"/slow-3s/" — handler sleeps 3 seconds before responding (timeout test).
type Upstream struct {
	srv *httptest.Server
}

func NewUpstream() *Upstream {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep paths first so the timeout test sees the deadline before any
		// response is written. Honor request cancellation so when the gateway
		// times out and drops the connection, this goroutine doesn't hold up
		// the test's teardown for the full 3s.
		if strings.Contains(r.URL.Path, "/slow-3s/") {
			select {
			case <-time.After(3 * time.Second):
			case <-r.Context().Done():
				return
			}
		}

		body, _ := io.ReadAll(r.Body)
		echo := Echo{
			Method:  r.Method,
			Path:    r.URL.Path,
			Host:    r.Host,
			Headers: r.Header,
			Body:    string(body),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(echo)
	}))
	return &Upstream{srv: srv}
}

// HostPort returns the host:port the gateway should forward to.
func (u *Upstream) HostPort() string {
	parsed, _ := url.Parse(u.srv.URL)
	return parsed.Host
}

func (u *Upstream) Close() { u.srv.Close() }
