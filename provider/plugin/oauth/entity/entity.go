package entity

import (
	"fmt"
	"time"
)

type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (eo ErrorObject) Error() string {
	return fmt.Sprintf("Error: %s, Code: %s", eo.Message, eo.Code)
}

type Error struct {
	HttpStatus int
	Errors     []ErrorObject
}

func (e Error) Error() string {
	errorString := "Service error because of:\n"

	for _, err := range e.Errors {
		errorString = errorString + err.Error() + "\n"
	}

	return errorString
}

type OauthPlugin struct {
	Config struct {
		Database string `yaml:"database"`

		AccessTokenTimeoutRaw       string `yaml:"access_token_timeout"`
		AuthorizationCodeTimeoutRaw string `yaml:"authorization_code_timeout"`

		RefreshToken struct {
			Timeout string `yaml:"timeout"`
			Active  bool   `yaml:"active"`
		} `yaml:"refresh_token"`
	} `yaml:"config"`
}

func (o OauthPlugin) DatabaseInstance() string {
	return o.Config.Database
}

func (o OauthPlugin) AccessTokenTimeout() (time.Duration, error) {
	return time.ParseDuration(o.Config.AccessTokenTimeoutRaw)
}

func (o OauthPlugin) AuthorizationCodeTimeout() (time.Duration, error) {
	return time.ParseDuration(o.Config.AuthorizationCodeTimeoutRaw)
}

func (o OauthPlugin) RefreshTokenTimeout() (time.Duration, error) {
	return time.ParseDuration(o.Config.RefreshToken.Timeout)
}
