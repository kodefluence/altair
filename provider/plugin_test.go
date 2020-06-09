package provider_test

import (
	"errors"
	"testing"

	"github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/provider"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Run gracefully", func(t *testing.T) {
		appBearer := mock.NewMockAppBearer(mockCtrl)
		appConfig := mock.NewMockAppConfig(mockCtrl)
		dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
		pluginBearer := mock.NewMockPluginBearer(mockCtrl)

		appBearer.EXPECT().Config().Return(appConfig)
		appConfig.EXPECT().PluginExists("oauth").Return(false)

		assert.NotPanics(t, func() {
			assert.Nil(t, provider.Plugin(appBearer, dbBearer, pluginBearer))
		})
	})

	t.Run("Oauth plugin error", func(t *testing.T) {
		appBearer := mock.NewMockAppBearer(mockCtrl)
		appConfig := mock.NewMockAppConfig(mockCtrl)
		dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
		pluginBearer := mock.NewMockPluginBearer(mockCtrl)

		appBearer.EXPECT().Config().Return(appConfig)
		appConfig.EXPECT().PluginExists("oauth").Return(true)
		pluginBearer.EXPECT().CompilePlugin("oauth", gomock.Any()).Return(errors.New("Unexpected error"))

		assert.NotPanics(t, func() {
			assert.NotNil(t, provider.Plugin(appBearer, dbBearer, pluginBearer))
		})
	})
}
