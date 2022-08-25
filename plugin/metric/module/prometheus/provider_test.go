package prometheus_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/mock"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/http"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus/usecase"
)

func TestProvider(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	appBearer := mock.NewMockAppBearer(mockController)
	appBearer.EXPECT().SetMetricProvider(usecase.NewPrometheus())
	appBearer.EXPECT().InjectController(http.NewPrometheusController())
	prometheus.Provide(appBearer)
}
