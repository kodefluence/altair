package plugin

import "github.com/codefluence-x/altair/core"

type provider struct{}

func Provider() core.PluginProviderDispatcher {
	return &provider{}
}

func (*provider) Oauth() core.PluginProvider {
	return nil
}
