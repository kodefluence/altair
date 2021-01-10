package downstream_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	coreEntity "github.com/codefluence-x/altair/entity"
	coreMock "github.com/codefluence-x/altair/mock"
	"github.com/codefluence-x/altair/provider/plugin/oauth/downstream"
	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestApplicationValidation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Name", func(t *testing.T) {
		t.Run("Return application-validation-plugin", func(t *testing.T) {
			oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
			oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
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

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", strings.NewReader(reqBody))
					routePath := coreEntity.RouterPath{Auth: "oauth_application"}

					entityOauthApplication := entity.OauthApplication{
						ID:           1,
						ClientUID:    clientUID,
						ClientSecret: clientSecret,
					}

					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthApplicationModel.EXPECT().OneByUIDandSecret(c, clientUID, clientSecret).Return(entityOauthApplication, nil)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
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

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", strings.NewReader(reqBody))
					routePath := coreEntity.RouterPath{Auth: "none"}

					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthApplicationModel.EXPECT().OneByUIDandSecret(c, gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
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

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", strings.NewReader(reqBody))
					routePath := coreEntity.RouterPath{Auth: "oauth_application"}

					entityOauthApplication := entity.OauthApplication{}

					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthApplicationModel.EXPECT().OneByUIDandSecret(c, clientUID, clientSecret).Return(entityOauthApplication, sql.ErrNoRows)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
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

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", strings.NewReader(reqBody))
					routePath := coreEntity.RouterPath{Auth: "oauth_application"}

					entityOauthApplication := entity.OauthApplication{}

					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthApplicationModel.EXPECT().OneByUIDandSecret(c, clientUID, clientSecret).Return(entityOauthApplication, errors.New("database is not available"))

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
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

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", strings.NewReader(reqBody))
					routePath := coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthApplicationModel.EXPECT().OneByUIDandSecret(c, gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
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

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", strings.NewReader(reqBody))
					routePath := coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthApplicationModel.EXPECT().OneByUIDandSecret(c, gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})

			t.Run("Body is nil", func(t *testing.T) {
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

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", nil)
					routePath := coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthApplicationModel.EXPECT().OneByUIDandSecret(c, gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
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

					r, _ := http.NewRequest("GET", "https://github.com/codefluence-x/altair", strings.NewReader(reqBody))
					routePath := coreEntity.RouterPath{Auth: "oauth_application"}

					oauthApplicationModel := mock.NewMockOauthApplicationModel(mockCtrl)
					oauthApplicationModel.EXPECT().OneByUIDandSecret(c, gomock.Any(), gomock.Any()).Times(0)

					oauthPlugin := downstream.NewApplicationValidation(oauthApplicationModel)
					err := oauthPlugin.Intervene(c, r, routePath)

					assert.NotNil(t, err)
				})
			})
		})
	})
}
