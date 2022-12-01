package entity

import (
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type MYSQLDatabaseConfig struct {
	Database              string `yaml:"database"`
	Username              string `yaml:"username"`
	Password              string `yaml:"password"`
	Host                  string `yaml:"host"`
	Port                  string `yaml:"port"`
	ConnectionMaxLifetime string `yaml:"connection_max_lifetime"`
	MaxIddleConnection    string `yaml:"max_iddle_connection"`
	MaxOpenConnection     string `yaml:"max_open_connection"`
}

func (m MYSQLDatabaseConfig) Driver() string {
	return "mysql"
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

func (m MYSQLDatabaseConfig) Dump() string {
	encodedContent, _ := yaml.Marshal(m)
	return string(encodedContent)
}
