package loader_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/loader"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseBearer(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dbConfig := entity.MYSQLDatabaseConfig{}

	databases := map[string]*sql.DB{
		"main_database": db,
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
					assert.Equal(t, db, loadedSQLDB)
					assert.Equal(t, dbConfig, loadedDBConfig)
				})
			})

			t.Run("Database is not found", func(t *testing.T) {
				t.Run("Return error", func(t *testing.T) {
					dbName := "this_is_not_exists_databases"

					loadedSQLDB, loadedDBConfig, err := dbBearer.Database(dbName)

					assert.Equal(t, loader.DatabasesIsNotExistsError, err)
					assert.Nil(t, loadedSQLDB)
					assert.Nil(t, loadedDBConfig)
				})
			})
		})
	})
}
