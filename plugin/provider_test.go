package plugin_test

import (
	"testing"

	"github.com/codefluence-x/altair/plugin"
	"github.com/stretchr/testify/assert"
)

func TestPluginProviderDispatcher(t *testing.T) {
	pluginProviderDispatcher := plugin.Provider()

	assert.NotPanics(t, func() {
		pluginProviderDispatcher.Oauth()
	})
}
