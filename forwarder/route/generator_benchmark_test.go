package route_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/forwarder/route"
	"github.com/codefluence-x/altair/provider/metric"
	"github.com/codefluence-x/altair/testhelper"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func BenchmarkRoute(b *testing.B) {
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

	err := route.Generator().Generate(gatewayEngine, metric.NewPrometheusMetric(), routeObjects, []core.DownStreamPlugin{})
	assert.Nil(b, err)

	srvTarget := &http.Server{
		Addr:    ":5002",
		Handler: targetEngine,
	}

	go func() {
		_ = srvTarget.ListenAndServe()
	}()

	// Given sleep time so the server can boot first
	time.Sleep(time.Millisecond * 50)

	for n := 0; n < b.N; n++ {
		assert.NotPanics(b, func() {
			rec := testhelper.PerformRequest(gatewayEngine, "GET", "/users/me", nil)
			assert.Equal(b, http.StatusOK, rec.Result().StatusCode)
		})
	}

	_ = srvTarget.Close()
}
