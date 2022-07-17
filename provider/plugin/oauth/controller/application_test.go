package controller_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/provider/plugin/oauth/controller"
	"github.com/kodefluence/altair/provider/plugin/oauth/mock"
	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Dispatch", func(t *testing.T) {
		applicationManager := mock.NewMockApplicationManager(mockCtrl)

		assert.NotPanics(t, func() {
			controller.NewApplication().List(applicationManager)
			controller.NewApplication().One(applicationManager)
			controller.NewApplication().Create(applicationManager)
			controller.NewApplication().Update(applicationManager)
		})
	})
}
