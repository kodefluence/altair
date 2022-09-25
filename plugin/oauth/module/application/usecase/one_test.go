package usecase_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application/usecase"
	"github.com/kodefluence/altair/plugin/oauth/module/application/usecase/mock"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)
	t.Run("One", func(t *testing.T) {
		t.Run("Given context and oauth application id", func(t *testing.T) {
			t.Run("Return oauth application data", func(t *testing.T) {
				oauthApplication := entity.OauthApplication{
					ID: 1,
				}

				ktx := kontext.Fabricate()
				formatterUsecase := newFormatter()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				oauthApplicationRepository.EXPECT().One(ktx, oauthApplication.ID, sqldb).Return(oauthApplication, nil)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
				oauthApplicationJSON, err := applicationManager.One(ktx, oauthApplication.ID)
				assert.Nil(t, err)
				assert.Equal(t, formatterUsecase.Application(oauthApplication), oauthApplicationJSON)
			})

			t.Run("Oauth application is not found", func(t *testing.T) {
				t.Run("Return 404", func(t *testing.T) {
					ktx := kontext.Fabricate()
					formatterUsecase := newFormatter()
					apierrorUsecase := apierror.Provide()
					oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

					oauthApplicationRepository.EXPECT().One(ktx, 1, sqldb).Return(entity.OauthApplication{}, exception.Throw(errors.New("not found"), exception.WithType(exception.NotFound)))

					applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
					_, err := applicationManager.One(ktx, 1)
					assert.NotNil(t, err)
					assert.Equal(t, http.StatusNotFound, err.HTTPStatus())
				})
			})

			t.Run("Unexpected error", func(t *testing.T) {
				t.Run("Return internal server error", func(t *testing.T) {
					ktx := kontext.Fabricate()
					formatterUsecase := newFormatter()
					apierrorUsecase := apierror.Provide()
					oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

					oauthApplicationRepository.EXPECT().One(ktx, 1, sqldb).Return(entity.OauthApplication{}, exception.Throw(errors.New("unexpected error"), exception.WithType(exception.Unexpected)))

					applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
					_, err := applicationManager.One(ktx, 1)
					assert.NotNil(t, err)
					assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus())
				})
			})
		})
	})
}
