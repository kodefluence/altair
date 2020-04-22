package plugin_test

import (
	"testing"

	"github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/plugin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	oauthAccessTokenModel := mock.NewMockOauthAccessTokenModel(mockCtrl)

	assert.NotPanics(t, func() {
		plugin.DownStream().Oauth(oauthAccessTokenModel)
	})
}
