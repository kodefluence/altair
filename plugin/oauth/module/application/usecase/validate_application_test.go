package usecase_test

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application/usecase"
	"github.com/kodefluence/altair/plugin/oauth/module/application/usecase/mock"
	"github.com/kodefluence/altair/util"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/stretchr/testify/assert"
)

func TestValidateApplication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)

	t.Run("ValidateApplication", func(t *testing.T) {
		t.Run("Given context and oauth application json data", func(t *testing.T) {
			t.Run("Return nil", func(t *testing.T) {
				data := entity.OauthApplicationJSON{
					OwnerID:     util.ValueToPointer(1),
					OwnerType:   util.ValueToPointer("confidential"),
					Description: util.ValueToPointer("This is description"),
					Scopes:      util.ValueToPointer("public users"),
				}

				formatterUsecase := newFormatter()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
				assert.Nil(t, applicationManager.ValidateApplication(data))
			})
		})

		t.Run("Given context and oauth application json data with empty owner_type", func(t *testing.T) {
			t.Run("Return validation error", func(t *testing.T) {
				data := entity.OauthApplicationJSON{
					OwnerID:     util.ValueToPointer(1),
					OwnerType:   nil,
					Description: util.ValueToPointer("This is description"),
					Scopes:      util.ValueToPointer("public users"),
				}

				formatterUsecase := newFormatter()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
				err := applicationManager.ValidateApplication(data)

				assert.NotNil(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, err.HTTPStatus())
				assert.Equal(t, "JSONAPI Error:\n[Validation error] Detail: Validation error because of: object `owner_type` is nil or not exists, Code: ERR1442\n", err.Error())
			})
		})

		t.Run("Given context and oauth application json data with invalid owner_type", func(t *testing.T) {
			t.Run("Return validation error", func(t *testing.T) {
				data := entity.OauthApplicationJSON{
					OwnerID:     util.ValueToPointer(1),
					OwnerType:   util.ValueToPointer("external"),
					Description: util.ValueToPointer("This is description"),
					Scopes:      util.ValueToPointer("public users"),
				}

				formatterUsecase := newFormatter()
				apierrorUsecase := apierror.Provide()
				oauthApplicationRepository := mock.NewMockOauthApplicationRepository(mockCtrl)

				applicationManager := usecase.NewApplicationManager(sqldb, oauthApplicationRepository, apierrorUsecase, formatterUsecase)
				err := applicationManager.ValidateApplication(data)

				assert.NotNil(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, err.HTTPStatus())
				assert.Equal(t, "JSONAPI Error:\n[Validation error] Detail: Validation error because of: object `owner_type` must be either of `confidential` or `public`, Code: ERR1442\n", err.Error())
			})
		})
	})
}
