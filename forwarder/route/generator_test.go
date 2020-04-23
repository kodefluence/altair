package route_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/forwarder/route"
	"github.com/codefluence-x/altair/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
						Path: map[string]struct{}{
							"/me":          struct{}{},
							"/details/:id": struct{}{},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				var downStreamPlugin []core.DownStreamPlugin
				downStreamPlugin = append(downStreamPlugin)

				err := route.Generator().Generate(gatewayEngine, routeObjects, downStreamPlugin)
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
					rec := mock.PerformRequest(gatewayEngine, "GET", "/users/me", nil)
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
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
						Path: map[string]struct{}{
							"/me":          struct{}{},
							"/details/:id": struct{}{},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				var downStreamPlugin []core.DownStreamPlugin

				oauthPlugin := mock.NewMockDownStreamPlugin(mockCtrl)
				oauthPlugin.EXPECT().Intervene(gomock.Any(), gomock.Any()).Return(nil)
				oauthPlugin.EXPECT().Name().AnyTimes().Return("oauth-plugin")

				downStreamPlugin = append(downStreamPlugin, oauthPlugin)

				err := route.Generator().Generate(gatewayEngine, routeObjects, downStreamPlugin)
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
					rec := mock.PerformRequest(gatewayEngine, "GET", "/users/me", nil)
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
						Path: map[string]struct{}{
							"/me":          struct{}{},
							"/details/:id": struct{}{},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				var downStreamPlugin []core.DownStreamPlugin

				oauthPlugin := mock.NewMockDownStreamPlugin(mockCtrl)
				oauthPlugin.EXPECT().Intervene(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
				oauthPlugin.EXPECT().Name().AnyTimes().Return("oauth-plugin")

				downStreamPlugin = append(downStreamPlugin, oauthPlugin)

				err := route.Generator().Generate(gatewayEngine, routeObjects, downStreamPlugin)
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
					rec := mock.PerformRequest(gatewayEngine, "GET", "/users/me", nil)
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
						Path: map[string]struct{}{
							"/me":          struct{}{},
							"/details/:id": struct{}{},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				var downStreamPlugin []core.DownStreamPlugin
				downStreamPlugin = append(downStreamPlugin)

				err := route.Generator().Generate(gatewayEngine, routeObjects, downStreamPlugin)
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
					rec := mock.PerformRequest(gatewayEngine, "GET", "/users/me/gusta", nil)
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
						Path: map[string]struct{}{
							"/me":          struct{}{},
							"/details/:id": struct{}{},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "POST", r)
				}

				var downStreamPlugin []core.DownStreamPlugin
				downStreamPlugin = append(downStreamPlugin)

				err := route.Generator().Generate(gatewayEngine, routeObjects, downStreamPlugin)
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
					rec := mock.PerformRequest(gatewayEngine, "POST", "/users/me", strings.NewReader(`{"id": 1, "state": "preparing"}`))
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
						Host:   "localhost:5002",
						Name:   "users",
						Prefix: "/users",
						Path: map[string]struct{}{
							"/me":  struct{}{},
							"/:id": struct{}{},
						},
					},
				)

				var downStreamPlugin []core.DownStreamPlugin
				downStreamPlugin = append(downStreamPlugin)

				err := route.Generator().Generate(gatewayEngine, routeObjects, downStreamPlugin)
				assert.NotNil(t, err)
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
						Path: map[string]struct{}{
							"/me":          struct{}{},
							"/details/:id": struct{}{},
						},
					},
				)

				for _, r := range routeObjects {
					buildTargetEngine(targetEngine, "GET", r)
				}

				var downStreamPlugin []core.DownStreamPlugin
				downStreamPlugin = append(downStreamPlugin)

				err := route.Generator().Generate(gatewayEngine, routeObjects, downStreamPlugin)
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
					rec := mock.PerformRequest(gatewayEngine, "GET", "/users/me", nil, func(req *http.Request) {
						req.Header.Add("foo", "bar")
					})
					assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
				})

				_ = srvTarget.Close()
			})
		})

		t.Run("Forwarding error", func(t *testing.T) {
			t.Run("Read request body error", func(t *testing.T) {

			})

			t.Run("Create new request error", func(t *testing.T) {

			})

			t.Run("Do the request error", func(t *testing.T) {

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
