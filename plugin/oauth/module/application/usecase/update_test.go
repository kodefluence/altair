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
	"github.com/kodefluence/altair/util"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)

	t.Run("Update", func(t *testing.T) {
		t.Run("Given context, id and oauth application update json", func(t *testing.T) {
			t.Run("When update process is success and find one process is success", func(t *testing.T) {
				t.Run("Then it will return oauth application json", func(t *testing.T) {
					oauthApplication := entity.OauthApplication{
						ID: 1,
					}

					ktx := kontext.Fabricate()
					formatterUsecase := newFormatter()
					apierrorUsecase := apierror.Provide()
					oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

					data := entity.OauthApplicationUpdateJSON{
						Description: util.ValueToPointer("New description"),
						Scopes:      util.ValueToPointer("users public"),
					}

					gomock.InOrder(
						oauthApplicationRepository.EXPECT().Update(gomock.Any(), oauthApplication.ID, entity.OauthApplicationUpdateable{
							Description: data.Description,
							Scopes:      data.Scopes,
						}, gomock.Any()).Return(nil),
						oauthApplicationRepository.EXPECT().One(ktx, 1, sqldb).Return(oauthApplication, nil),
					)

					applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
					oauthApplicationJSON, err := applicationManager.Update(ktx, oauthApplication.ID, data)
					assert.Nil(t, err)
					assert.Equal(t, formatterUsecase.Application(oauthApplication), oauthApplicationJSON)
				})
			})

			t.Run("When update process failed", func(t *testing.T) {
				t.Run("Then it will return error", func(t *testing.T) {
					oauthApplication := entity.OauthApplication{
						ID: 1,
					}

					ktx := kontext.Fabricate()
					formatterUsecase := newFormatter()
					apierrorUsecase := apierror.Provide()
					oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

					data := entity.OauthApplicationUpdateJSON{
						Description: util.ValueToPointer("New description"),
						Scopes:      util.ValueToPointer("users public"),
					}

					gomock.InOrder(
						oauthApplicationRepository.EXPECT().Update(gomock.Any(), oauthApplication.ID, entity.OauthApplicationUpdateable{
							Description: data.Description,
							Scopes:      data.Scopes,
						}, gomock.Any()).Return(exception.Throw(errors.New("unexpected"))),
					)

					applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
					_, err := applicationManager.Update(ktx, oauthApplication.ID, data)
					assert.NotNil(t, err)
					assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus())
				})
			})
		})
	})
}
