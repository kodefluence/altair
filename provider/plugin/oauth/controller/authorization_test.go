package controller_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/provider/plugin/oauth/controller"
	"github.com/kodefluence/altair/provider/plugin/oauth/mock"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	authService := mock.NewMockAuthorization(mockCtrl)

	t.Run("Dispatch", func(t *testing.T) {
		assert.NotPanics(t, func() {
			controller.NewAuthorization().Grant(authService)
			controller.NewAuthorization().Revoke(authService)
			controller.NewAuthorization().Token(authService)
		})
	})
}
