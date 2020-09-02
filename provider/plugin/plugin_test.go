package plugin_test

import (
	"testing"

	"github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/provider/plugin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			err := plugin.Oauth(appBearer, dbBearer, pluginBearer)
			require.NoError(t, err)
		})
	})
}
