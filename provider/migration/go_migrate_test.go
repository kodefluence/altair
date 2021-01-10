package migration_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/provider/migration"
	"github.com/codefluence-x/altair/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestGoMigrate(t *testing.T) {

	t.Run("Migrator", func(t *testing.T) {
		t.Run("Run gracefully", func(t *testing.T) {
			t.Run("Return migrator", func(t *testing.T) {
				db, mockDB, err := sqlmock.New()
				if err != nil {
					panic(err)
				}

				testhelper.GenerateTempTestFiles("./migration/", "some content", "01_migration.up.sql", 0666)

				dbConfig := entity.MYSQLDatabaseConfig{
					MigrationSource: "file://migration",
					Database:        "altair_development",
				}

				mockDB.ExpectQuery(`SELECT GET_LOCK\(\?, 10\)`).WillReturnRows(sqlmock.NewRows([]string{
					"GET_LOCK(1, 10)",
				}).AddRow(1))

				mockDB.ExpectQuery(`SHOW TABLES LIKE "db_versions"`).WillReturnRows(sqlmock.NewRows([]string{
					"TABLES IN DATABASES",
				}).AddRow("db_versions"))

				mockDB.ExpectExec(`SELECT RELEASE_LOCK\(\?\)`).WillReturnResult(sqlmock.NewResult(0, 1))

				migrationProvider := migration.GoMigrate(db, dbConfig)
				migrator, err := migrationProvider.Migrator()

				assert.Nil(t, err)
				assert.NotNil(t, migrator)

				testhelper.RemoveTempTestFiles("./migration/")
			})
		})

		t.Run("Database instantiate failed", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				db, mockDB, err := sqlmock.New()
				if err != nil {
					panic(err)
				}

				dbConfig := entity.MYSQLDatabaseConfig{
					MigrationSource: "file://instantiate_failed",
					Database:        "altair_development",
				}

				mockDB.ExpectQuery(`SELECT GET_LOCK\(\?, 10\)`).WillReturnError(errors.New("Unexpected error"))

				migrationProvider := migration.GoMigrate(db, dbConfig)
				migrator, err := migrationProvider.Migrator()

				assert.NotNil(t, err)
				assert.Nil(t, migrator)
			})
		})

		t.Run("Migrator instantiate failed", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				db, mockDB, err := sqlmock.New()
				if err != nil {
					panic(err)
				}

				dbConfig := entity.MYSQLDatabaseConfig{
					MigrationSource: "file://this_path_should_not_exists",
					Database:        "altair_development",
				}

				mockDB.ExpectQuery(`SELECT GET_LOCK\(\?, 10\)`).WillReturnRows(sqlmock.NewRows([]string{
					"GET_LOCK(1, 10)",
				}).AddRow(1))

				mockDB.ExpectQuery(`SHOW TABLES LIKE "db_versions"`).WillReturnRows(sqlmock.NewRows([]string{
					"TABLES IN DATABASES",
				}).AddRow("db_versions"))

				mockDB.ExpectExec(`SELECT RELEASE_LOCK\(\?\)`).WillReturnResult(sqlmock.NewResult(0, 1))

				migrationProvider := migration.GoMigrate(db, dbConfig)
				migrator, err := migrationProvider.Migrator()

				assert.NotNil(t, err)
				assert.Nil(t, migrator)
			})
		})
	})
}
