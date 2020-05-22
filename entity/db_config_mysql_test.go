package entity_test

import (
	"testing"
	"time"

	"github.com/codefluence-x/altair/entity"
	"github.com/stretchr/testify/assert"
)

func TestMYSQLDatabaseConfig(t *testing.T) {
	MYSQLConfig := entity.MYSQLDatabaseConfig{
		Database:              "altair_development",
		Username:              "some_username",
		Password:              "some_password",
		Host:                  "localhost",
		Port:                  "3306",
		ConnectionMaxLifetime: "120s",
		MaxIddleConnection:    "100",
		MaxOpenConnection:     "100",
		MigrationSource:       "file://migration",
	}

	expectedDatabase := "altair_development"
	expectedUsername := "some_username"
	expectedPassword := "some_password"
	expectedHost := "localhost"
	expectedPort := 3306
	expectedConnMaxLifetime := time.Second * 120
	expectedMaxIddleConn := 100
	expectedMaxOpenConn := 100
	expectedMigrationSource := "file://migration"

	assert.Equal(t, "mysql", MYSQLConfig.Driver())
	assert.Equal(t, expectedDatabase, MYSQLConfig.DBDatabase())
	assert.Equal(t, expectedUsername, MYSQLConfig.DBUsername())
	assert.Equal(t, expectedPassword, MYSQLConfig.DBPassword())
	assert.Equal(t, expectedMigrationSource, MYSQLConfig.DBMigrationSource())
	assert.Equal(t, expectedHost, MYSQLConfig.DBHost())

	actualPort, err := MYSQLConfig.DBPort()
	assert.Nil(t, err)
	assert.Equal(t, expectedPort, actualPort)

	actualConnMaxLifetime, err := MYSQLConfig.DBConnectionMaxLifetime()
	assert.Nil(t, err)
	assert.Equal(t, expectedConnMaxLifetime, actualConnMaxLifetime)

	actualMaxIddleConn, err := MYSQLConfig.DBMaxIddleConn()
	assert.Nil(t, err)
	assert.Equal(t, expectedMaxIddleConn, actualMaxIddleConn)

	actualMaxOpenConn, err := MYSQLConfig.DBMaxOpenConn()
	assert.Nil(t, err)
	assert.Equal(t, expectedMaxOpenConn, actualMaxOpenConn)
}
