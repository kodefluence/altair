package forwarder_test

import (
	"testing"

	"github.com/codefluence-x/altair/forwarder"
	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {

	assert.NotPanics(t, func() {
		forwarder.Route().Generator()
		forwarder.Route().Compiler()
	})
}
