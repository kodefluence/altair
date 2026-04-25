package cfg

import (
	"errors"
	"fmt"

	"github.com/kodefluence/monorepo/db"

	"github.com/kodefluence/altair/core"
)

// ErrDatabasesIsNotExists thrown when database is not exists or not initialized yet in database bearer
var ErrDatabasesIsNotExists = errors.New("Database is not exists")

type databaseBearer struct {
	databases map[string]db.DB
	configs   map[string]core.DatabaseConfig
}

// DatabaseBearer handling on retrieval database instance
func DatabaseBearer(databases map[string]db.DB, configs map[string]core.DatabaseConfig) core.DatabaseBearer {
	return &databaseBearer{databases: databases, configs: configs}
}

func (d *databaseBearer) Database(dbName string) (db.DB, core.DatabaseConfig, error) {
	sqldb, ok := d.databases[dbName]
	if !ok {
		return nil, nil, ErrDatabasesIsNotExists
	}

	config, ok := d.configs[dbName]
	if !ok {
		return nil, nil, fmt.Errorf("database `%s` is connected but has no matching config entry", dbName)
	}

	return sqldb, config, nil
}
