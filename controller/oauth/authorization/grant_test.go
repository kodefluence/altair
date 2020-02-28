package authorization_test

import (
	"testing"

	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGrant(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Method", func(t *testing.T) {
		authorizationService := mock.NewMockAuthorization(mockCtrl)
		assert.Equal(t, "POST", controller.Oauth().Authorization().Grant(authorizationService).Method())
	})

	t.Run("Path", func(t *testing.T) {
		authorizationService := mock.NewMockAuthorization(mockCtrl)
		assert.Equal(t, "/oauth/authorizations", controller.Oauth().Authorization().Grant(authorizationService).Path())
	})

	t.Run("Control", func(t *testing.T) {
		t.Run("Given request with json body", func(t *testing.T) {

		})
	})
}
