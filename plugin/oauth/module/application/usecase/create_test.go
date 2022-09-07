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
	"github.com/kodefluence/altair/plugin/oauth/module/formatter"
	"github.com/kodefluence/monorepo/db"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)

	t.Run("Given context and oauth application json data", func(t *testing.T) {
		t.Run("Return last inserted id", func(t *testing.T) {
			ktx := kontext.Fabricate()
			formatterUsecase := formatter.Provide()
			apierrorUsecase := apierror.Provide()
			oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

			applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)

			oauthApplication := entity.OauthApplication{
				ID:        1,
				OwnerType: "confidential",
			}
			oauthApplicationJSON := formatterUsecase.Application(oauthApplication)

			gomock.InOrder(
				oauthApplicationRepository.EXPECT().Create(ktx, gomock.Any(), sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthApplicationInsertable, tx db.TX) (int, exception.Exception) {
					assert.Equal(t, oauthApplication.OwnerType, data.OwnerType)
					return oauthApplication.ID, nil
				}),
				oauthApplicationRepository.EXPECT().One(ktx, oauthApplication.ID, sqldb).Return(oauthApplication, nil),
			)

			result, err := applicationManager.Create(ktx, oauthApplicationJSON)

			assert.Nil(t, err)
			assert.Equal(t, oauthApplicationJSON, result)
		})

		t.Run("Validation error", func(t *testing.T) {
			t.Run("Return unprocessable entity", func(t *testing.T) {
				ktx := kontext.Fabricate()
				formatterUsecase := formatter.Provide()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)

				oauthApplication := entity.OauthApplication{
					ID:        1,
					OwnerType: "",
				}
				oauthApplicationJSON := formatterUsecase.Application(oauthApplication)

				_, err := applicationManager.Create(ktx, oauthApplicationJSON)

				assert.NotNil(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, err.HTTPStatus())
			})
		})

		t.Run("Unexpected error", func(t *testing.T) {
			t.Run("Return internal server error", func(t *testing.T) {
				ktx := kontext.Fabricate()
				formatterUsecase := formatter.Provide()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)

				oauthApplication := entity.OauthApplication{
					ID:        1,
					OwnerType: "confidential",
				}
				oauthApplicationJSON := formatterUsecase.Application(oauthApplication)

				gomock.InOrder(
					oauthApplicationRepository.EXPECT().Create(ktx, gomock.Any(), sqldb).DoAndReturn(func(ktx kontext.Context, data entity.OauthApplicationInsertable, tx db.TX) (int, exception.Exception) {
						assert.Equal(t, oauthApplication.OwnerType, data.OwnerType)
						return 0, exception.Throw(errors.New("unexpected errors"))
					}),
				)

				_, err := applicationManager.Create(ktx, oauthApplicationJSON)

				assert.NotNil(t, err)
				assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus())
			})
		})
	})
}
