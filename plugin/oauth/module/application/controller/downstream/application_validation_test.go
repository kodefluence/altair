package downstream_test

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	coreEntity "github.com/kodefluence/altair/entity"
	coreMock "github.com/kodefluence/altair/mock"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/downstream"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/downstream/mock"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/exception"
	"github.com/stretchr/testify/assert"
)

func TestApplicationValidation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sqldb := mockdb.NewMockDB(mockCtrl)

	t.Run("Name", func(t *testing.T) {
		t.Run("Return application-validation-plugin", func(t *testing.T) {
			oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
			oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
			assert.Equal(t, "application-validation-plugin", oauthPlugin.Name())
		})
	})

	t.Run("Intervene", func(t *testing.T) {
		t.Run("Given gin.Context and http.Request", func(t *testing.T) {
			t.Run("Normal scenario", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					clientUID := "application_uid"
					clientSecret := "client_secret"

					reqBody := fmt.Sprintf(`{"client_uid":"%s","client_secret":"%s","username":"altair","password":"handsomeeagle"}`, clientUID, clientSecret)

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", strings.NewReader(reqBody))
					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					entityOauthApplication := entity.OauthApplication{
						ID:           1,
						ClientUID:    clientUID,
						ClientSecret: clientSecret,
					}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), clientUID, clientSecret, sqldb).Return(entityOauthApplication, nil)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.Nil(t, err)
				})
			})

			t.Run("Auth is not oauth application", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					clientUID := "application_uid"
					clientSecret := "client_secret"

					reqBody := fmt.Sprintf(`{"client_uid":"%s","client_secret":"%s","username":"altair","password":"handsomeeagle"}`, clientUID, clientSecret)

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", strings.NewReader(reqBody))
					routePath := &coreEntity.RouterPath{Auth: "none"}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.Nil(t, err)
				})
			})

			t.Run("Oauth application is not found", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusUnauthorized).AnyTimes()

					c.Writer = responseWritterMock

					clientUID := "application_uid"
					clientSecret := "client_secret"

					reqBody := fmt.Sprintf(`{"client_uid":"%s","client_secret":"%s","username":"altair","password":"handsomeeagle"}`, clientUID, clientSecret)

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", strings.NewReader(reqBody))
					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					entityOauthApplication := entity.OauthApplication{}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), clientUID, clientSecret, sqldb).Return(entityOauthApplication, exception.Throw(sql.ErrNoRows, exception.WithType(exception.NotFound)))

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})

			t.Run("Database is currently not available", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusServiceUnavailable).AnyTimes()

					c.Writer = responseWritterMock

					clientUID := "application_uid"
					clientSecret := "client_secret"

					reqBody := fmt.Sprintf(`{"client_uid":"%s","client_secret":"%s","username":"altair","password":"handsomeeagle"}`, clientUID, clientSecret)

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", strings.NewReader(reqBody))
					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					entityOauthApplication := entity.OauthApplication{}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), clientUID, clientSecret, sqldb).Return(entityOauthApplication, exception.Throw(errors.New("data is not available")))

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})

			t.Run("Client uid is not provided", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusUnprocessableEntity).AnyTimes()

					c.Writer = responseWritterMock

					clientSecret := "client_secret"

					reqBody := fmt.Sprintf(`{"client_secret":"%s","username":"altair","password":"handsomeeagle"}`, clientSecret)

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", strings.NewReader(reqBody))
					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})

			t.Run("Client secret is not provided", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusUnprocessableEntity).AnyTimes()

					c.Writer = responseWritterMock

					clientUID := "client_uid"

					reqBody := fmt.Sprintf(`{"client_uid":"%s","username":"altair","password":"handsomeeagle"}`, clientUID)

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", strings.NewReader(reqBody))
					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})

			t.Run("Body is nil", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusUnprocessableEntity).AnyTimes()

					c.Writer = responseWritterMock

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)
					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})

			t.Run("Body is nil but have header with CLIENT_UID and CLIENT_SECRET", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {
					clientUID := "application_uid"
					clientSecret := "client_secret"

					header := http.Header{}
					header.Add("CLIENT_UID", clientUID)
					header.Add("CLIENT_SECRET", clientSecret)

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: header,
					}

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", nil)
					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					entityOauthApplication := entity.OauthApplication{
						ID:           1,
						ClientUID:    clientUID,
						ClientSecret: clientSecret,
					}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), clientUID, clientSecret, sqldb).Return(entityOauthApplication, nil)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)
					assert.Nil(t, err)
				})
			})

			t.Run("GetBody returning an error", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusServiceUnavailable).AnyTimes()

					c.Writer = responseWritterMock

					clientUID := "application_uid"
					clientSecret := "client_secret"

					reqBody := fmt.Sprintf(`{"client_uid":"%s","client_secret":"%s","username":"altair","password":"handsomeeagle"}`, clientUID, clientSecret)
					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", strings.NewReader(reqBody))
					r.GetBody = func() (io.ReadCloser, error) {
						return nil, errors.New("unexpected error")
					}

					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})

			t.Run("Body is not json", func(t *testing.T) {
				t.Run("Return nil", func(t *testing.T) {

					c := &gin.Context{}
					c.Request = &http.Request{
						Header: http.Header{},
					}

					responseWritterMock := coreMock.NewMockResponseWriter(mockCtrl)
					responseWritterMock.EXPECT().WriteHeaderNow().AnyTimes()
					responseWritterMock.EXPECT().WriteHeader(gomock.Any()).AnyTimes()
					responseWritterMock.EXPECT().Status().Return(http.StatusUnprocessableEntity).AnyTimes()

					c.Writer = responseWritterMock

					clientUID := "application_uid"
					clientSecret := "client_secret"

					reqBody := fmt.Sprintf(`client_uid=%s client_secrent=%s`, clientUID, clientSecret)

					r, _ := http.NewRequest("GET", "https://github.com/kodefluence/altair", strings.NewReader(reqBody))
					routePath := &coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationRepo := mock.NewMockOauthApplicationRepository(mockCtrl)
					oauthApplicationRepo.EXPECT().OneByUIDandSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationRepo, sqldb)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})
		})
	})
}
