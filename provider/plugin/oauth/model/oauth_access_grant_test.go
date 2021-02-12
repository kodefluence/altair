package model_test

import (
	"errors"
	"testing"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/model"
	"github.com/codefluence-x/altair/provider/plugin/oauth/query"
	mockdb "github.com/codefluence-x/monorepo/db/mock"
	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOauthAccessGrant(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("One", func(t *testing.T) {
		t.Run("Given context and oauth access grant ID", func(t *testing.T) {
			t.Run("When database operation complete it will return oauth access grant data", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				row := mockdb.NewMockRow(mockCtrl)

				expectedData := entity.OauthAccessGrant{
					ID:   1,
					Code: "code",
				}

				sqldb.EXPECT().QueryRowContext(gomock.Any(), "oauth-access-grant-one", query.SelectOneOauthAccessGrant, expectedData.ID).Return(row)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) exception.Exception {
					val0, _ := dest[0].(*int)
					*val0 = expectedData.ID

					val1, _ := dest[4].(*string)
					*val1 = expectedData.Code
					return nil
				})

				oauthAccessGrant := model.NewOauthAccessGrant()
				data, err := oauthAccessGrant.One(kontext.Fabricate(), expectedData.ID, sqldb)

				assert.Nil(t, err)
				assert.Equal(t, expectedData, data)
			})
		})
	})

	t.Run("OneByCode", func(t *testing.T) {
		t.Run("Given context and oauth access grant ID", func(t *testing.T) {
			t.Run("When database operation complete it will return oauth access grant data", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				row := mockdb.NewMockRow(mockCtrl)

				expectedData := entity.OauthAccessGrant{
					ID:   1,
					Code: "code",
				}

				sqldb.EXPECT().QueryRowContext(gomock.Any(), "oauth-access-grant-one-by-code", query.SelectOneOauthAccessGrantByCode, expectedData.Code).Return(row)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) exception.Exception {
					val0, _ := dest[0].(*int)
					*val0 = expectedData.ID

					val1, _ := dest[4].(*string)
					*val1 = expectedData.Code
					return nil
				})

				oauthAccessGrant := model.NewOauthAccessGrant()
				data, err := oauthAccessGrant.OneByCode(kontext.Fabricate(), expectedData.Code, sqldb)

				assert.Nil(t, err)
				assert.Equal(t, expectedData, data)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Given context and access grant insertable", func(t *testing.T) {
			t.Run("When database operation complete it will create oauth access grant data", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				result := mockdb.NewMockResult(mockCtrl)

				expectedID := 1
				insertable := entity.OauthAccessGrantInsertable{
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               "token",
					Scopes:             "public users stores",
					RedirectURI:        "https://github.com",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-access-grant-create", query.InsertOauthAccessGrant,
					insertable.OauthApplicationID, insertable.ResourceOwnerID, insertable.Scopes, insertable.Code, insertable.RedirectURI, insertable.ExpiresIn,
				).Return(result, nil)

				result.EXPECT().LastInsertId().Return(int64(expectedID), nil)

				oauthAccessGrantModel := model.NewOauthAccessGrant()
				ID, err := oauthAccessGrantModel.Create(kontext.Fabricate(), insertable, sqldb)

				assert.Equal(t, expectedID, ID)
				assert.Nil(t, err)
			})

			t.Run("When database operation failed then it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)

				insertable := entity.OauthAccessGrantInsertable{
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               "token",
					Scopes:             "public users stores",
					RedirectURI:        "https://github.com",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-access-grant-create", query.InsertOauthAccessGrant,
					insertable.OauthApplicationID, insertable.ResourceOwnerID, insertable.Scopes, insertable.Code, insertable.RedirectURI, insertable.ExpiresIn,
				).Return(nil, exception.Throw(errors.New("unexpected")))

				oauthAccessGrantModel := model.NewOauthAccessGrant()
				_, err := oauthAccessGrantModel.Create(kontext.Fabricate(), insertable, sqldb)

				assert.Equal(t, exception.Unexpected, err.Type())
				assert.NotNil(t, err)
			})

			t.Run("When database operation is complete but get last inserted id error then it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				result := mockdb.NewMockResult(mockCtrl)

				insertable := entity.OauthAccessGrantInsertable{
					OauthApplicationID: 1,
					ResourceOwnerID:    1,
					Code:               "token",
					Scopes:             "public users stores",
					RedirectURI:        "https://github.com",
					ExpiresIn:          time.Now().Add(time.Hour * 24),
				}

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-access-grant-create", query.InsertOauthAccessGrant,
					insertable.OauthApplicationID, insertable.ResourceOwnerID, insertable.Scopes, insertable.Code, insertable.RedirectURI, insertable.ExpiresIn,
				).Return(result, nil)

				result.EXPECT().LastInsertId().Return(int64(0), exception.Throw(errors.New("unexpected error")))

				oauthAccessGrantModel := model.NewOauthAccessGrant()
				ID, err := oauthAccessGrantModel.Create(kontext.Fabricate(), insertable, sqldb)

				assert.Equal(t, 0, ID)
				assert.NotNil(t, err)
				assert.Equal(t, exception.Unexpected, err.Type())
			})
		})
	})

	t.Run("Revoke", func(t *testing.T) {
		t.Run("Given context and authorization code", func(t *testing.T) {
			t.Run("When database operation complete it will return nil", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				result := mockdb.NewMockResult(mockCtrl)

				code := "code"

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-access-grant-revoke", query.RevokeAuthorizationCode, code).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)

				oauthAccessGrantModel := model.NewOauthAccessGrant()
				err := oauthAccessGrantModel.Revoke(kontext.Fabricate(), code, sqldb)
				assert.Nil(t, err)
			})

			t.Run("When database operation complete but there is no updated rows it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				result := mockdb.NewMockResult(mockCtrl)

				code := "code"

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-access-grant-revoke", query.RevokeAuthorizationCode, code).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(0), nil)

				oauthAccessGrantModel := model.NewOauthAccessGrant()
				err := oauthAccessGrantModel.Revoke(kontext.Fabricate(), code, sqldb)
				assert.NotNil(t, err)
				assert.Equal(t, exception.NotFound, err.Type())
			})

			t.Run("When database operation failed it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				code := "code"

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-access-grant-revoke", query.RevokeAuthorizationCode, code).Return(nil, exception.Throw(errors.New("unexpected error")))

				oauthAccessGrantModel := model.NewOauthAccessGrant()
				err := oauthAccessGrantModel.Revoke(kontext.Fabricate(), code, sqldb)
				assert.NotNil(t, err)
			})

			t.Run("When database operation complete but get rows affected error it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				result := mockdb.NewMockResult(mockCtrl)

				code := "code"

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-access-grant-revoke", query.RevokeAuthorizationCode, code).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(0), exception.Throw(errors.New("unexpected error")))

				oauthAccessGrantModel := model.NewOauthAccessGrant()
				err := oauthAccessGrantModel.Revoke(kontext.Fabricate(), code, sqldb)
				assert.NotNil(t, err)
			})
		})
	})
}
