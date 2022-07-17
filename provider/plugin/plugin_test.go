package plugin_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/mock"
	"github.com/kodefluence/altair/provider/plugin"
	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Oauth", func(t *testing.T) {
		appBearer := mock.NewMockAppBearer(mockCtrl)
		appConfig := mock.NewMockAppConfig(mockCtrl)
		dbBearer := mock.NewMockDatabaseBearer(mockCtrl)
		pluginBearer := mock.NewMockPluginBearer(mockCtrl)

		appBearer.EXPECT().Config().Return(appConfig)
		appConfig.EXPECT().PluginExists("oauth").Return(false)

		assert.NotPanics(t, func() {
			plugin.Oauth(appBearer, dbBearer, pluginBearer)
		})
	})
}
