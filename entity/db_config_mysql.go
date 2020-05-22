package entity

import (
	"strconv"
	"time"
)

type MYSQLDatabaseConfig struct {
	Database              string
	Username              string
	Password              string
	Host                  string
	Port                  string
	ConnectionMaxLifetime string
	MaxIddleConnection    string
	MaxOpenConnection     string
	MigrationSource       string
}

func (m MYSQLDatabaseConfig) Driver() string {
	return "mysql"
}

func (m MYSQLDatabaseConfig) DBMigrationSource() string {
	return m.MigrationSource
}

func (m MYSQLDatabaseConfig) DBHost() string {
	return m.Host
}

func (m MYSQLDatabaseConfig) DBPort() (int, error) {
	return strconv.Atoi(m.Port)
}

func (m MYSQLDatabaseConfig) DBUsername() string {
	return m.Username
}

func (m MYSQLDatabaseConfig) DBPassword() string {
	return m.Password
}

func (m MYSQLDatabaseConfig) DBDatabase() string {
	return m.Database
}

func (m MYSQLDatabaseConfig) DBConnectionMaxLifetime() (time.Duration, error) {
	return time.ParseDuration(m.ConnectionMaxLifetime)
}

func (m MYSQLDatabaseConfig) DBMaxIddleConn() (int, error) {
	return strconv.Atoi(m.MaxIddleConnection)
}

func (m MYSQLDatabaseConfig) DBMaxOpenConn() (int, error) {
	return strconv.Atoi(m.MaxOpenConnection)
}
