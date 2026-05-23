package usecase_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/mock"
	"github.com/kodefluence/altair/module/router/usecase"
	"github.com/kodefluence/altair/testhelper"
)

func TestGenerator(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gin.SetMode(gin.ReleaseMode)

	t.Run("Generate", func(t *testing.T) {
		t.Run("Call target services routes", func(t *testing.T) {
			t.Run("Run gracefully", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:5002",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":          {Auth: "none"},
							"/details/:id": {Auth: "none"},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				var downStreamController []module.DownstreamController

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":5002",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/me", nil)
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})

			t.Run("Return original status code from the target server", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:9081",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":          {Auth: "none"},
							"/details/:id": {Auth: "none"},
						},
					},
				)

				targetEngine.GET("/users/me", func(c *gin.Context) {
					c.Status(http.StatusInternalServerError)
				})

				var downStreamController []module.DownstreamController

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":9081",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/me", nil)
					assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})

			t.Run("Return original status code from the target server", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:9701",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":          {Auth: "none"},
							"/details/:id": {Auth: "none"},
						},
					},
				)

				targetEngine.GET("/users/me", func(c *gin.Context) {
					assert.Equal(t, "bar", c.Query("foo"))
					assert.Equal(t, "case", c.Query("cool"))
					c.Status(http.StatusInternalServerError)
				})

				var downStreamController []module.DownstreamController

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":9701",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/me?foo=bar&cool=case", nil)
					assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})
		})

		t.Run("Call target services routes with downstream plugins", func(t *testing.T) {
			t.Run("Run gracefully", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:5011",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":          {Auth: "none"},
							"/details/:id": {Auth: "none"},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				oauthPlugin := mock.NewMockDownstreamController(mockCtrl)
				oauthPlugin.EXPECT().Intervene(gomock.Any(), gomock.Any(), routeObjects[0].Path["/me"]).Return(nil)
				oauthPlugin.EXPECT().Name().AnyTimes().Return("oauth-plugin")

				var downStreamController []module.DownstreamController
				downStreamController = append(downStreamController, oauthPlugin)

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":5011",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/me", nil)
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})

			t.Run("Run gracefully with configured route path", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:5012",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":            {},
							"/authorization": {Auth: "oauth"},
							"/details/:id":   {},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				oauthPlugin := mock.NewMockDownstreamController(mockCtrl)
				oauthPlugin.EXPECT().Intervene(gomock.Any(), gomock.Any(), routeObjects[0].Path["/authorization"]).Return(nil)
				oauthPlugin.EXPECT().Name().AnyTimes().Return("oauth-plugin")

				var downStreamController []module.DownstreamController
				downStreamController = append(downStreamController, oauthPlugin)

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":5012",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/authorization", nil)
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})

			t.Run("Run gracefully with wildcard route path", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "oauth",
						Host:   "localhost:5013",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/details/:id": {Auth: "oauth"},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				oauthPlugin := mock.NewMockDownstreamController(mockCtrl)
				oauthPlugin.EXPECT().Intervene(gomock.Any(), gomock.Any(), routeObjects[0].Path["/details/:id"]).Return(nil)
				oauthPlugin.EXPECT().Name().AnyTimes().Return("oauth-plugin")

				var downStreamController []module.DownstreamController
				downStreamController = append(downStreamController, oauthPlugin)

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":5013",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/details/me", nil)
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})

			t.Run("Downstream plugins error", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:5011",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":          {Auth: "none"},
							"/details/:id": {Auth: "none"},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				oauthPlugin := mock.NewMockDownstreamController(mockCtrl)
				oauthPlugin.EXPECT().Intervene(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
				oauthPlugin.EXPECT().Name().AnyTimes().Return("oauth-plugin")

				var downStreamController []module.DownstreamController
				downStreamController = append(downStreamController, oauthPlugin)

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":5011",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/me", nil)
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})
		})

		t.Run("Call target services routes and the routes is not found", func(t *testing.T) {
			t.Run("Run gracefully and return 404 status", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:5003",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":          {Auth: "none"},
							"/details/:id": {Auth: "none"},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				var downStreamController []module.DownstreamController
				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":5003",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/me/gusta", nil)
					assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})
		})

		t.Run("Call target services with post and body", func(t *testing.T) {
			t.Run("Run gracefully", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:5004",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":          {Auth: "none"},
							"/details/:id": {Auth: "none"},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "POST", r)
				}

				var downStreamController []module.DownstreamController

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":5004",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "POST", "/users/me", strings.NewReader(`{"id": 1, "state": "preparing"}`))
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})
		})

		t.Run("Route wildcard error", func(t *testing.T) {
			t.Run("Return error when generate the routing", func(t *testing.T) {
				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:5019",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":  {Auth: "none"},
							"/:id": {Auth: "none"},
						},
					},
				)

				var downStreamController []module.DownstreamController

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)
			})
		})

		t.Run("Call target services routes with custom headers", func(t *testing.T) {
			t.Run("Run gracefully", func(t *testing.T) {
				targetEngine := gin.Default()

				gatewayEngine := gin.New()

				var routeObjects []entity.RouteObject
				routeObjects = append(
					routeObjects,
					entity.RouteObject{
						Auth:   "none",
						Host:   "localhost:5005",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]entity.RouterPath{
							"/me":          {Auth: "none"},
							"/details/:id": {Auth: "none"},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				var downStreamController []module.DownstreamController

				err := usecase.NewGenerator(downStreamController, []module.MetricController{testhelper.NewDummyMetric()}).Generate(gatewayEngine, routeObjects)
				assert.Nil(t, err)

				srvTarget := &http.Server{
					Addr:    ":5005",
					Handler: targetEngine,
				}

				go func() {
					_ = srvTarget.ListenAndServe()
				}()

				// Given sleep time so the server can boot first
				time.Sleep(time.Millisecond * 100)

				assert.NotPanics(t, func() {
					rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/me", nil, func(req *http.Request) {
						req.Header.Add("foo", "bar")
					})
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})
		})

		t.Run("Rejects request body larger than configured cap with 413", func(t *testing.T) {
			// Pin the assumption: WithMaxRequestBodySize(N) makes the
			// gateway return 413 (Request Entity Too Large) for any client
			// request whose body exceeds N. Default zero-value cap means
			// unlimited and is exercised by the existing POST/body tests.
			gatewayEngine := gin.New()

			routeObjects := []entity.RouteObject{{
				Auth:   "none",
				Host:   "127.0.0.1:5096",
				Name:   "limit",
				Prefix: "/limit",
				Path:   map[string]entity.RouterPath{"/upload": {Auth: "none"}},
			}}

			generator := usecase.NewGenerator(
				nil,
				[]module.MetricController{testhelper.NewDummyMetric()},
				usecase.WithMaxRequestBodySize(8), // 8 bytes
			)
			assert.Nil(t, generator.Generate(gatewayEngine, routeObjects))

			// Upstream should never be reached — but if it is, fail loud.
			upstreamHit := false
			targetEngine := gin.New()
			targetEngine.POST("/limit/upload", func(c *gin.Context) {
				upstreamHit = true
				c.Status(http.StatusOK)
			})

			srvTarget := &http.Server{Addr: ":5096", Handler: targetEngine}
			go func() { _ = srvTarget.ListenAndServe() }()
			defer srvTarget.Close()
			time.Sleep(100 * time.Millisecond)

			rec := testhelper.PerformRequest(gatewayEngine, "POST", "/limit/upload", strings.NewReader("0123456789ABCDEF"))
			assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Result().StatusCode)
			assert.False(t, upstreamHit, "oversized request must not reach upstream")
		})

		t.Run("Forwards configured proxy host to upstream", func(t *testing.T) {
			// Pin the assumption: WithProxyHost replaces the per-request
			// os.Getenv("PROXY_HOST") read with a value captured once at
			// construction. The upstream sees that exact value as the Host
			// header — regardless of what the env var holds at request time.
			gatewayEngine := gin.New()

			routeObjects := []entity.RouteObject{{
				Auth:   "none",
				Host:   "127.0.0.1:5097",
				Name:   "hostcheck",
				Prefix: "/h",
				Path:   map[string]entity.RouterPath{"/probe": {Auth: "none"}},
			}}

			generator := usecase.NewGenerator(
				nil,
				[]module.MetricController{testhelper.NewDummyMetric()},
				usecase.WithProxyHost("captured.example.com"),
			)
			assert.Nil(t, generator.Generate(gatewayEngine, routeObjects))

			gotHost := make(chan string, 1)
			targetEngine := gin.New()
			targetEngine.GET("/h/probe", func(c *gin.Context) {
				gotHost <- c.Request.Host
				c.Status(http.StatusOK)
			})

			srvTarget := &http.Server{Addr: ":5097", Handler: targetEngine}
			go func() { _ = srvTarget.ListenAndServe() }()
			defer srvTarget.Close()
			time.Sleep(100 * time.Millisecond)

			rec := testhelper.PerformRequest(gatewayEngine, "GET", "/h/probe", nil)
			assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

			select {
			case got := <-gotHost:
				assert.Equal(t, "captured.example.com", got)
			case <-time.After(2 * time.Second):
				t.Fatal("upstream never received the request")
			}
		})

		t.Run("Streams large response body intact", func(t *testing.T) {
			// Pin the assumption: a multi-MB upstream response is delivered
			// to the client byte-for-byte, with the upstream status code
			// preserved. Pre-fix, generator.go did io.ReadAll(proxyRes.Body)
			// then c.Writer.Write(resp) — buffering the entire body in RAM.
			// Post-fix uses io.Copy and must still yield a byte-equal payload.
			gatewayEngine := gin.New()

			routeObjects := []entity.RouteObject{{
				Auth:   "none",
				Host:   "127.0.0.1:5098",
				Name:   "big",
				Prefix: "/big",
				Path:   map[string]entity.RouterPath{"/payload": {Auth: "none"}},
			}}

			// 4MB so we'd notice if buffering is reintroduced; small enough
			// that the test stays under a second.
			payload := make([]byte, 4*1024*1024)
			for i := range payload {
				payload[i] = byte(i % 251)
			}

			generator := usecase.NewGenerator(
				nil,
				[]module.MetricController{testhelper.NewDummyMetric()},
			)
			assert.Nil(t, generator.Generate(gatewayEngine, routeObjects))

			targetEngine := gin.New()
			targetEngine.GET("/big/payload", func(c *gin.Context) {
				c.Data(http.StatusTeapot, "application/octet-stream", payload)
			})

			srvTarget := &http.Server{Addr: ":5098", Handler: targetEngine}
			go func() { _ = srvTarget.ListenAndServe() }()
			defer srvTarget.Close()
			time.Sleep(100 * time.Millisecond)

			rec := testhelper.PerformRequest(gatewayEngine, "GET", "/big/payload", nil)
			assert.Equal(t, http.StatusTeapot, rec.Result().StatusCode)
			assert.Equal(t, len(payload), rec.Body.Len())
			assert.Equal(t, payload, rec.Body.Bytes())
		})

		t.Run("Upstream exceeds configured timeout returns 502 promptly", func(t *testing.T) {
			// Pin the assumption: when an upstream hangs longer than the
			// configured client timeout, the gateway must return 502 (the
			// existing failure mode for client.Do errors) rather than waiting
			// indefinitely. Without WithUpstreamTimeout the proxy used a
			// zero-value http.Client{} that blocked forever — see
			// generator.go:173 prior to this change.
			gatewayEngine := gin.New()

			routeObjects := []entity.RouteObject{{
				Auth:   "none",
				Host:   "127.0.0.1:5099",
				Name:   "slow",
				Prefix: "/slow",
				Path: map[string]entity.RouterPath{
					"/hang": {Auth: "none"},
				},
			}}

			generator := usecase.NewGenerator(
				nil,
				[]module.MetricController{testhelper.NewDummyMetric()},
				usecase.WithUpstreamTimeout(100*time.Millisecond),
			)
			assert.Nil(t, generator.Generate(gatewayEngine, routeObjects))

			targetEngine := gin.New()
			targetEngine.GET("/slow/hang", func(c *gin.Context) {
				time.Sleep(3 * time.Second)
				c.Status(http.StatusOK)
			})

			srvTarget := &http.Server{Addr: ":5099", Handler: targetEngine}
			go func() { _ = srvTarget.ListenAndServe() }()
			defer srvTarget.Close()
			time.Sleep(100 * time.Millisecond)

			start := time.Now()
			rec := testhelper.PerformRequest(gatewayEngine, "GET", "/slow/hang", nil)
			elapsed := time.Since(start)

			assert.Equal(t, http.StatusBadGateway, rec.Result().StatusCode)
			// Must complete well before the upstream's 3s sleep — otherwise
			// the timeout isn't being honored.
			assert.Less(t, elapsed, 1*time.Second, "expected timeout to fire well before upstream sleep finishes, got %s", elapsed)
		})

		// TODO: add test
		t.Run("Forwarding error", func(t *testing.T) {
			t.Run("Do the request error", func(t *testing.T) {

			})

			t.Run("Read request body error", func(t *testing.T) {

			})

			t.Run("Create new request error", func(t *testing.T) {

			})

			t.Run("Read response body error", func(t *testing.T) {

			})

			t.Run("Write response body error", func(t *testing.T) {

			})
		})
	})
}

func buildTargetEngine(targetEngine *gin.Engine, method string, routeObject entity.RouteObject) {
	for p := range routeObject.Path {
		targetEngine.Handle(method, fmt.Sprintf("%s%s", routeObject.Prefix, p), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}
}
