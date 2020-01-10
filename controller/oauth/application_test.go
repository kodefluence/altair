package oauth_test

import (
	"testing"

	"github.com/codefluence-x/altair/controller/oauth"
	"github.com/codefluence-x/altair/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Dispatch", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)

		assert.NotPanics(t, func() {
			oauth.Application().List(applicationManager)
			oauth.Application().One(applicationManager)
		})
	})
}
