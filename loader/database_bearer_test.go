package loader_test

import (
	"testing"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/monorepo/db"
	mockdb "github.com/codefluence-x/monorepo/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseBearer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbConfig := entity.MYSQLDatabaseConfig{}

	sqldb := mockdb.NewMockDB(mockCtrl)
	databases := map[string]db.DB{
		"main_database": sqldb,
	}

	configs := map[string]core.DatabaseConfig{
		"main_database": dbConfig,
	}

	dbBearer := loader.DatabaseBearer(databases, configs)

	t.Run("Database", func(t *testing.T) {
		t.Run("Given database name", func(t *testing.T) {
			t.Run("Database is found", func(t *testing.T) {
				t.Run("Return database instance", func(t *testing.T) {
					dbName := "main_database"

					loadedSQLDB, loadedDBConfig, err := dbBearer.Database(dbName)

					assert.Nil(t, err)
					assert.Equal(t, sqldb, loadedSQLDB)
					assert.Equal(t, dbConfig, loadedDBConfig)
				})
			})

			t.Run("Database is not found", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					dbName := "this_is_not_exists_databases"

					loadedSQLDB, loadedDBConfig, err := dbBearer.Database(dbName)

					assert.Equal(t, loader.ErrDatabasesIsNotExists, err)
					assert.Nil(t, loadedSQLDB)
					assert.Nil(t, loadedDBConfig)
				})
			})
		})
	})
}
