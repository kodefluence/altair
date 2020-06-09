package controller_test

import (
	"testing"

	"github.com/codefluence-x/altair/provider/plugin/oauth/controller"
	"github.com/codefluence-x/altair/provider/plugin/oauth/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Dispatch", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)

		assert.NotPanics(t, func() {
			controller.Application().List(applicationManager)
			controller.Application().One(applicationManager)
			controller.Application().Create(applicationManager)
		})
	})
}
