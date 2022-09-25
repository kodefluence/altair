package entity

import "time"

// OauthPlugin holds all config variables
type OauthPlugin struct {
	Config PluginConfig `yaml:"config"`
}

// PluginConfig holds all config variables for oauth plugin
type PluginConfig struct {
	Database string `yaml:"database"`

	AccessTokenTimeoutRaw       string `yaml:"access_token_timeout"`
	AuthorizationCodeTimeoutRaw string `yaml:"authorization_code_timeout"`

	RefreshToken struct {
		Timeout string `yaml:"timeout"`
		Active  bool   `yaml:"active"`
	} `yaml:"refresh_token"`

	ImplicitGrant struct {
		Active bool `yaml:"active"`
	} `yaml:"implicit_grant"`
}

type RefreshTokenConfig struct {
	Timeout time.Duration
	Active  bool
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
