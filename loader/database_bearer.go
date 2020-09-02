package loader

import (
	"database/sql"

	"github.com/codefluence-x/altair/core"
	"github.com/pkg/errors"
)

var DatabasesIsNotExistsError = errors.New("Database is not exists")

type databaseBearer struct {
	databases map[string]*sql.DB
	configs   map[string]core.DatabaseConfig
}

func DatabaseBearer(databases map[string]*sql.DB, configs map[string]core.DatabaseConfig) core.DatabaseBearer {
	return &databaseBearer{databases: databases, configs: configs}
}

func (d *databaseBearer) Database(dbName string) (*sql.DB, core.DatabaseConfig, error) {
	db, ok := d.databases[dbName]
	if !ok {
		return nil, nil, DatabasesIsNotExistsError
	}

	return db, d.configs[dbName], nil
}
