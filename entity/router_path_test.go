package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/altair/entity"
)

func TestRouterPath(t *testing.T) {
	routerPath := entity.RouterPath{
		Auth:  "oauth",
		Scope: "public",
	}

	assert.Equal(t, routerPath.Auth, routerPath.GetAuth())
	assert.Equal(t, routerPath.Scope, routerPath.GetScope())
}
