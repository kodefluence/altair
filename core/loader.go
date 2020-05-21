package core

import (
	"time"
)

type DatabaseLoader interface {
	Compile(configPath string) (map[string]DatabaseConfig, error)
}

type DatabaseConfig interface {
	Driver() string
	DBHost() string
	DBPort() (int, error)
	DBUsername() string
	DBPassword() string
	DBDatabase() string
	DBConnectionMaxLifetime() (time.Duration, error)
	DBMaxIddleConn() (int, error)
	DBMaxOpenConn() (int, error)
}
