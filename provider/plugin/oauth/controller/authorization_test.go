package controller_test

import (
	"testing"

	"github.com/codefluence-x/altair/provider/plugin/oauth/controller"
	"github.com/codefluence-x/altair/provider/plugin/oauth/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	authService := mock.NewMockAuthorization(mockCtrl)

	t.Run("Dispatch", func(t *testing.T) {
		assert.NotPanics(t, func() {
			controller.Authorization().Grant(authService)
			controller.Authorization().Revoke(authService)
			controller.Authorization().Token(authService)
		})
	})
}
