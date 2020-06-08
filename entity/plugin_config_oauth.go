package entity

import "time"

type OauthPlugin struct {
	Config struct {
		Database string `yaml:"database"`

		AccessTokenTimeoutRaw       string `yaml:"access_token_timeout"`
		AuthorizationCodeTimeoutRaw string `yaml:"authorization_code_timeout"`
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
