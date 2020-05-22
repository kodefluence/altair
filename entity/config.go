package entity

import "time"

type Plugin struct {
	Plugin string `yaml:"plugin"`
	Raw    []byte
}

type OauthPlugins struct {
	Config struct {
		Database string `yaml:"database"`

		AccessTokenTimeoutRaw string        `yaml:"access_token_timeout"`
		AccessTokenTimeout    time.Duration `yaml:"-"`

		AuthorizationCodeTimeoutRaw string        `yaml:"authorization_code_timeout"`
		AuthorizationCodeTimeout    time.Duration `yaml:"-"`
	} `yaml:"config"`
}
