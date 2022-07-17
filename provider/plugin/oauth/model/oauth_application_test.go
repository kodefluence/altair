package model_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/model"
	"github.com/kodefluence/altair/provider/plugin/oauth/query"
	mockdb "github.com/kodefluence/monorepo/db/mock"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestOauthApplication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Paginate", func(t *testing.T) {
		t.Run("Given offset and limit", func(t *testing.T) {
			t.Run("When database operation complete then it will return array of oauth applications data", func(t *testing.T) {
				offset := 0
				limit := 10

				sqldb := mockdb.NewMockDB(mockCtrl)
				rows := mockdb.NewMockRows(mockCtrl)

				expectedOauthApplications := []entity.OauthApplication{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
				}

				sqldb.EXPECT().QueryContext(gomock.Any(), "oauth-application-paginate", query.PaginateOauthApplication, offset, limit).Return(rows, nil)
				rows.EXPECT().Next().Return(true).Times(len(expectedOauthApplications))

				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(),
				).DoAndReturn(func(dest ...interface{}) exception.Exception {
					val, _ := dest[0].(*int)
					*val = expectedOauthApplications[0].ID
					fmt.Println(val)
					return nil
				})
				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(),
				).DoAndReturn(func(dest ...interface{}) exception.Exception {
					val, _ := dest[0].(*int)
					*val = expectedOauthApplications[1].ID
					fmt.Println(val)
					return nil
				})

				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				rows.EXPECT().Close()

				oauthApplicationModel := model.NewOauthApplication()
				oauthApplications, err := oauthApplicationModel.Paginate(kontext.Fabricate(), offset, limit, sqldb)

				assert.Nil(t, err)
				assert.Equal(t, expectedOauthApplications, oauthApplications)
			})

			t.Run("When database operation failed then it will return error", func(t *testing.T) {
				offset := 0
				limit := 10

				sqldb := mockdb.NewMockDB(mockCtrl)
				sqldb.EXPECT().QueryContext(gomock.Any(), "oauth-application-paginate", query.PaginateOauthApplication, offset, limit).Return(nil, exception.Throw(errors.New("unexpected")))

				oauthApplicationModel := model.NewOauthApplication()
				oauthApplications, err := oauthApplicationModel.Paginate(kontext.Fabricate(), offset, limit, sqldb)

				assert.NotNil(t, err)
				assert.Equal(t, exception.Unexpected, err.Type())
				assert.Equal(t, 0, len(oauthApplications))
			})

			t.Run("When database operation complete but row scan error it will return error", func(t *testing.T) {
				offset := 0
				limit := 10

				sqldb := mockdb.NewMockDB(mockCtrl)
				rows := mockdb.NewMockRows(mockCtrl)

				expectedOauthApplications := []entity.OauthApplication{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
				}

				sqldb.EXPECT().QueryContext(gomock.Any(), "oauth-application-paginate", query.PaginateOauthApplication, offset, limit).Return(rows, nil)
				rows.EXPECT().Next().Return(true).Times(len(expectedOauthApplications))

				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(),
				).DoAndReturn(func(dest ...interface{}) exception.Exception {
					val, _ := dest[0].(*int)
					*val = expectedOauthApplications[0].ID
					fmt.Println(val)
					return nil
				})
				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(),
				).DoAndReturn(func(dest ...interface{}) exception.Exception {
					return exception.Throw(errors.New("unexpected"))
				})

				rows.EXPECT().Close()

				oauthApplicationModel := model.NewOauthApplication()
				oauthApplications, err := oauthApplicationModel.Paginate(kontext.Fabricate(), offset, limit, sqldb)

				assert.NotNil(t, err)
				assert.Equal(t, exception.Unexpected, err.Type())
				assert.Equal(t, 1, len(oauthApplications))
			})
		})
	})

	t.Run("Count", func(t *testing.T) {
		t.Run("When database operation complete it will return total data of oauth_applications", func(t *testing.T) {
			sqldb := mockdb.NewMockDB(mockCtrl)
			row := mockdb.NewMockRow(mockCtrl)

			expectedTotal := 100

			sqldb.EXPECT().QueryRowContext(gomock.Any(), "oauth-application-count", query.CountOauthApplication).Return(row)
			row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) exception.Exception {
				val, _ := dest[0].(*int)
				*val = expectedTotal
				return nil
			})

			oauthApplicationModel := model.NewOauthApplication()
			total, err := oauthApplicationModel.Count(kontext.Fabricate(), sqldb)

			assert.Nil(t, err)
			assert.Equal(t, expectedTotal, total)
		})

		t.Run("When database operation complete but row scan error then it will return exception", func(t *testing.T) {
			sqldb := mockdb.NewMockDB(mockCtrl)
			row := mockdb.NewMockRow(mockCtrl)

			expectedTotal := 0

			sqldb.EXPECT().QueryRowContext(gomock.Any(), "oauth-application-count", query.CountOauthApplication).Return(row)
			row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) exception.Exception {
				return exception.Throw(errors.New("unexpected error"))
			})

			oauthApplicationModel := model.NewOauthApplication()
			total, err := oauthApplicationModel.Count(kontext.Fabricate(), sqldb)

			assert.NotNil(t, err)
			assert.Equal(t, exception.Unexpected, err.Type())
			assert.Equal(t, expectedTotal, total)
		})
	})

	t.Run("One", func(t *testing.T) {
		t.Run("Given oauth application id", func(t *testing.T) {
			t.Run("When database operation complete then it will return oauth application data", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				row := mockdb.NewMockRow(mockCtrl)

				expectedData := entity.OauthApplication{
					ID:        1,
					OwnerType: "confidential",
				}

				sqldb.EXPECT().QueryRowContext(gomock.Any(), "oauth-application-one", query.SelectOneOauthApplication, expectedData.ID).Return(row)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) exception.Exception {
					val0, _ := dest[0].(*int)
					*val0 = expectedData.ID
					val2, _ := dest[2].(*string)
					*val2 = expectedData.OwnerType
					return nil
				})

				oauthApplicationModel := model.NewOauthApplication()
				data, err := oauthApplicationModel.One(kontext.Fabricate(), expectedData.ID, sqldb)

				assert.Nil(t, err)
				assert.Equal(t, expectedData, data)
			})

			t.Run("When database operation complete but scan failed then it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				row := mockdb.NewMockRow(mockCtrl)

				expectedData := entity.OauthApplication{}

				sqldb.EXPECT().QueryRowContext(gomock.Any(), "oauth-application-one", query.SelectOneOauthApplication, expectedData.ID).Return(row)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) exception.Exception {
					return exception.Throw(errors.New("unexpected error"))
				})

				oauthApplicationModel := model.NewOauthApplication()
				data, err := oauthApplicationModel.One(kontext.Fabricate(), expectedData.ID, sqldb)

				assert.NotNil(t, err)
				assert.Equal(t, exception.Unexpected, err.Type())
				assert.Equal(t, expectedData, data)
			})
		})
	})

	t.Run("OneByUIDandSecret", func(t *testing.T) {
		t.Run("Given client uid and client secret", func(t *testing.T) {
			t.Run("When database operation complete then it will return oauth application data", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				row := mockdb.NewMockRow(mockCtrl)

				expectedData := entity.OauthApplication{
					ID:           1,
					OwnerType:    "confidential",
					ClientUID:    "uid",
					ClientSecret: "secret",
				}

				sqldb.EXPECT().QueryRowContext(gomock.Any(), "oauth-application-one-by-id-and-secret", query.SelectOneByUIDandSecret, expectedData.ClientUID, expectedData.ClientSecret).Return(row)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) exception.Exception {
					val0, _ := dest[0].(*int)
					*val0 = expectedData.ID
					val2, _ := dest[2].(*string)
					*val2 = expectedData.OwnerType
					val3, _ := dest[5].(*string)
					*val3 = expectedData.ClientUID
					val4, _ := dest[6].(*string)
					*val4 = expectedData.ClientSecret
					return nil
				})

				oauthApplicationModel := model.NewOauthApplication()
				data, err := oauthApplicationModel.OneByUIDandSecret(kontext.Fabricate(), expectedData.ClientUID, expectedData.ClientSecret, sqldb)

				assert.Nil(t, err)
				assert.Equal(t, expectedData, data)
			})

			t.Run("When database operation complete but scan failed then it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				row := mockdb.NewMockRow(mockCtrl)

				expectedData := entity.OauthApplication{}

				sqldb.EXPECT().QueryRowContext(gomock.Any(), "oauth-application-one-by-id-and-secret", query.SelectOneByUIDandSecret, "uid", "secret").Return(row)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) exception.Exception {
					return exception.Throw(errors.New("unexpected error"))
				})

				oauthApplicationModel := model.NewOauthApplication()
				data, err := oauthApplicationModel.OneByUIDandSecret(kontext.Fabricate(), "uid", "secret", sqldb)

				assert.NotNil(t, err)
				assert.Equal(t, exception.Unexpected, err.Type())
				assert.Equal(t, expectedData, data)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Given context and oauth application json entity", func(t *testing.T) {
			t.Run("When database operation is complete then it will return last inserted id", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				result := mockdb.NewMockResult(mockCtrl)

				expectedID := 1
				insertable := entity.OauthApplicationInsertable{
					OwnerID:      1,
					OwnerType:    "confidential",
					Description:  "-",
					Scopes:       "public user",
					ClientUID:    "client-uid",
					ClientSecret: "client-secret",
				}

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-application-create", query.InsertOauthApplication,
					insertable.OwnerID,
					insertable.OwnerType,
					insertable.Description,
					insertable.Scopes,
					insertable.ClientUID,
					insertable.ClientSecret,
				).Return(result, nil)

				result.EXPECT().LastInsertId().Return(int64(expectedID), nil)

				oauthApplicationModel := model.NewOauthApplication()
				ID, err := oauthApplicationModel.Create(kontext.Fabricate(), insertable, sqldb)

				assert.Equal(t, expectedID, ID)
				assert.Nil(t, err)
			})

			t.Run("When database operation failed then it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)

				insertable := entity.OauthApplicationInsertable{
					OwnerID:      1,
					OwnerType:    "confidential",
					Description:  "-",
					Scopes:       "public user",
					ClientUID:    "client-uid",
					ClientSecret: "client-secret",
				}

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-application-create", query.InsertOauthApplication,
					insertable.OwnerID,
					insertable.OwnerType,
					insertable.Description,
					insertable.Scopes,
					insertable.ClientUID,
					insertable.ClientSecret,
				).Return(nil, exception.Throw(errors.New("unexpected")))

				oauthApplicationModel := model.NewOauthApplication()
				_, err := oauthApplicationModel.Create(kontext.Fabricate(), insertable, sqldb)

				assert.Equal(t, exception.Unexpected, err.Type())
				assert.NotNil(t, err)
			})

			t.Run("When database operation is complete but get last inserted id error then it will return error", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				result := mockdb.NewMockResult(mockCtrl)

				insertable := entity.OauthApplicationInsertable{
					OwnerID:      1,
					OwnerType:    "confidential",
					Description:  "-",
					Scopes:       "public user",
					ClientUID:    "client-uid",
					ClientSecret: "client-secret",
				}

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-application-create", query.InsertOauthApplication,
					insertable.OwnerID,
					insertable.OwnerType,
					insertable.Description,
					insertable.Scopes,
					insertable.ClientUID,
					insertable.ClientSecret,
				).Return(result, nil)

				result.EXPECT().LastInsertId().Return(int64(0), exception.Throw(errors.New("unexpected error")))

				oauthApplicationModel := model.NewOauthApplication()
				ID, err := oauthApplicationModel.Create(kontext.Fabricate(), insertable, sqldb)

				assert.Equal(t, 0, ID)
				assert.NotNil(t, err)
				assert.Equal(t, exception.Unexpected, err.Type())
			})
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("Given context, oauth application id and oauth application updateable", func(t *testing.T) {
			t.Run("When database operation is complete it will return value", func(t *testing.T) {
				sqldb := mockdb.NewMockDB(mockCtrl)
				result := mockdb.NewMockResult(mockCtrl)

				ID := 1
				data := entity.OauthApplicationUpdateable{
					Description: "X",
					Scopes:      "public user",
				}

				sqldb.EXPECT().ExecContext(gomock.Any(), "oauth-application-update", query.UpdateOauthApplication, data.Description, data.Scopes, ID).Return(result, nil)

				oauthApplicationModel := model.NewOauthApplication()
				err := oauthApplicationModel.Update(kontext.Fabricate(), ID, data, sqldb)
				assert.Nil(t, err)
			})
		})
	})
}
