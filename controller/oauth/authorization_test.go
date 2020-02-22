package oauth_test

import (
	"testing"

	"github.com/codefluence-x/altair/controller/oauth"
	"github.com/codefluence-x/altair/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	authService := mock.NewMockAuthorization(mockCtrl)

	t.Run("Dispatch", func(t *testing.T) {
		assert.NotPanics(t, func() {
			oauth.Authorization().Grant(authService)
		})
	})
}
