package usecase_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application/usecase"
	"github.com/kodefluence/altair/plugin/oauth/module/application/usecase/mock"

	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

func TestList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)

	t.Run("List", func(t *testing.T) {
		t.Run("Given limit and offset", func(t *testing.T) {
			t.Run("Return a formatted result", func(t *testing.T) {
				oauthApplications := []entity.OauthApplication{
					{ID: 1},
					{ID: 2},
				}

				ktx := kontext.Fabricate()
				formatterUsecase := newFormatter()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				gomock.InOrder(
					oauthApplicationRepository.EXPECT().Paginate(ktx, 0, 10, sqldb).Return(oauthApplications, nil),
					oauthApplicationRepository.EXPECT().Count(ktx, sqldb).Return(len(oauthApplications), nil),
				)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
				oauthApplicationJSON, total, err := applicationManager.List(ktx, 0, 10)
				assert.Nil(t, err)
				assert.Equal(t, len(oauthApplications), total)
				assert.Equal(t, formatterUsecase.ApplicationList(oauthApplications), oauthApplicationJSON)

			})

			t.Run("Error oauth application paginate return internal server error", func(t *testing.T) {
				ktx := kontext.Fabricate()
				formatterUsecase := newFormatter()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				gomock.InOrder(
					oauthApplicationRepository.EXPECT().Paginate(ktx, 0, 10, sqldb).Return([]entity.OauthApplication(nil), exception.Throw(errors.New("unexpected errors"))),
				)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
				oauthApplicationJSON, total, err := applicationManager.List(ktx, 0, 10)

				assert.NotNil(t, err)
				assert.Equal(t, err.HTTPStatus(), http.StatusInternalServerError)
				assert.Equal(t, 0, total)
				assert.Equal(t, []entity.OauthApplicationJSON(nil), oauthApplicationJSON)
			})

			t.Run("Error oauth application count return internal server error", func(t *testing.T) {
				oauthApplications := []entity.OauthApplication{
					{ID: 1},
					{ID: 2},
				}

				ktx := kontext.Fabricate()
				formatterUsecase := newFormatter()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				gomock.InOrder(
					oauthApplicationRepository.EXPECT().Paginate(ktx, 0, 10, sqldb).Return(oauthApplications, nil),
					oauthApplicationRepository.EXPECT().Count(ktx, sqldb).Return(0, exception.Throw(errors.New("unexpected errors"))),
				)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
				oauthApplicationJSON, total, err := applicationManager.List(ktx, 0, 10)

				assert.NotNil(t, err)
				assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus())
				assert.Equal(t, 0, total)
				assert.Equal(t, []entity.OauthApplicationJSON(nil), oauthApplicationJSON)
			})
		})
	})
}
