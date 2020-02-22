package controller_test

import (
	"testing"

	"github.com/codefluence-x/altair/controller"
	"github.com/stretchr/testify/assert"
)

func TestOauthDispatcher(t *testing.T) {

	t.Run("Dispatch", func(t *testing.T) {
		assert.NotPanics(t, func() {
			controller.Oauth().Application()
			controller.Oauth().Authorization()
		})
	})
}
