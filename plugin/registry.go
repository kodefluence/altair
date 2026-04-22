package plugin

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/metric"
	"github.com/kodefluence/altair/plugin/oauth"
)

// Registry returns every plugin compiled into this binary. This is the single
// audit surface for "what does this binary know how to do?" — adding a plugin
// is one line here plus a Plugin struct in the plugin's package.
//
// The order returned here is *input* order to the topological sorter in
// plugin.Load / plugin.LoadCommand; after the sort is applied it is
// deterministic regardless of slice order. We list alphabetically to keep
// diffs stable.
func Registry() []module.Plugin {
	return []module.Plugin{
		&metric.Plugin{},
		&oauth.Plugin{},
	}
}
