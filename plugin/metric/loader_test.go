package metric_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/mock"
	"github.com/kodefluence/altair/plugin/metric"
	"github.com/kodefluence/altair/plugin/metric/entity"
	dummyUsecase "github.com/kodefluence/altair/plugin/metric/module/dummy/usecase"
	promHttp "github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/http"
	promUsecase "github.com/kodefluence/altair/plugin/metric/module/prometheus/usecase"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {

	t.Run("When plugin metric is not set, then it will return dummy metric", func(t *testing.T) {
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		appBearer := mock.NewMockAppBearer(mockController)
		appConfig := mock.NewMockAppConfig(mockController)
		pluginBearer := mock.NewMockPluginBearer(mockController)

		gomock.InOrder(
			appBearer.EXPECT().Config().Return(appConfig),
			appConfig.EXPECT().PluginExists("metric").Return(false),
			appBearer.EXPECT().SetMetricProvider(dummyUsecase.NewDummy()),
		)

		assert.Nil(t, metric.Load(appBearer, pluginBearer))
	})

	t.Run("When plugin metric is set with prometheus, then it will return prometheus metric", func(t *testing.T) {
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		appBearer := mock.NewMockAppBearer(mockController)
		appConfig := mock.NewMockAppConfig(mockController)
		pluginBearer := mock.NewMockPluginBearer(mockController)

		gomock.InOrder(
			appBearer.EXPECT().Config().Return(appConfig),
			appConfig.EXPECT().PluginExists("metric").Return(true),
			pluginBearer.EXPECT().CompilePlugin("metric", gomock.Any()).DoAndReturn(func(pluginName string, injectedStruct interface{}) error {
				v, _ := injectedStruct.(*entity.MetricPlugin)
				v.Config.Provider = "prometheus"
				return nil
			}),
			appBearer.EXPECT().SetMetricProvider(promUsecase.NewPrometheus()),
			appBearer.EXPECT().InjectController(promHttp.NewPrometheusController()),
		)

		assert.Nil(t, metric.Load(appBearer, pluginBearer))
	})

	t.Run("When plugin metric is set with prometheus but the plugin bearer return error when compiling, then it will return error", func(t *testing.T) {
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		appBearer := mock.NewMockAppBearer(mockController)
		appConfig := mock.NewMockAppConfig(mockController)
		pluginBearer := mock.NewMockPluginBearer(mockController)

		expectedError := errors.New("expectedError")
		gomock.InOrder(
			appBearer.EXPECT().Config().Return(appConfig),
			appConfig.EXPECT().PluginExists("metric").Return(true),
			pluginBearer.EXPECT().CompilePlugin("metric", gomock.Any()).DoAndReturn(func(pluginName string, injectedStruct interface{}) error {
				return expectedError
			}),
		)

		assert.Equal(t, expectedError, metric.Load(appBearer, pluginBearer))
	})

	t.Run("When plugin metric is set with unsupported provider, then it will return error", func(t *testing.T) {
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		appBearer := mock.NewMockAppBearer(mockController)
		appConfig := mock.NewMockAppConfig(mockController)
		pluginBearer := mock.NewMockPluginBearer(mockController)

		gomock.InOrder(
			appBearer.EXPECT().Config().Return(appConfig),
			appConfig.EXPECT().PluginExists("metric").Return(true),
			pluginBearer.EXPECT().CompilePlugin("metric", gomock.Any()).DoAndReturn(func(pluginName string, injectedStruct interface{}) error {
				v, _ := injectedStruct.(*entity.MetricPlugin)
				v.Config.Provider = "datadog"
				return nil
			}),
		)

		assert.Equal(t, fmt.Errorf("Metric plugin `%s` is currently not supported", "datadog"), metric.Load(appBearer, pluginBearer))
	})

}