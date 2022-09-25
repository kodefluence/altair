package usecase

import "time"

type Formatter struct {
	tokenExpiresIn        time.Duration
	codeExpiresIn         time.Duration
	refreshTokenExpiresIn time.Duration
}

func NewFormatter(tokenExpiresIn time.Duration, codeExpiresIn time.Duration, refreshTokenExpiresIn time.Duration) *Formatter {
	return &Formatter{
		tokenExpiresIn:        tokenExpiresIn,
		codeExpiresIn:         codeExpiresIn,
		refreshTokenExpiresIn: refreshTokenExpiresIn,
	}
}
