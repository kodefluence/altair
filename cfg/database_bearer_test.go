package cfg_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/cfg"
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/entity"
	"github.com/kodefluence/monorepo/db"
	mockdb "github.com/kodefluence/monorepo/db/mock"
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

	dbBearer := cfg.DatabaseBearer(databases, configs)

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

					assert.Equal(t, cfg.ErrDatabasesIsNotExists, err)
					assert.Nil(t, loadedSQLDB)
					assert.Nil(t, loadedDBConfig)
				})
			})
		})
	})
}
