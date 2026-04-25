package plugin

import (
	"fmt"
	"sort"

	"github.com/kodefluence/monorepo/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
)

// Load builds a PluginContext for every active plugin in Registry() and calls
// Plugin.Load in topologically-sorted order. A plugin is "active" iff it is
// listed in app.yml `plugins:` AND has a matching config/plugin/<name>.yml
// (AND-gate).
func Load(appBearer core.AppBearer, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer, apiError module.ApiError, appModule module.App) error {
	return run(Registry(), appBearer, pluginBearer, dbBearer, apiError, appModule, func(p module.Plugin, ctx module.PluginContext) error {
		return p.Load(ctx)
	})
}

// LoadCommand is the `altair plugin ...` equivalent of Load: it iterates the
// same sorted registry but calls Plugin.LoadCommand, which registers CLI
// subcommands instead of HTTP/downstream/metric controllers.
func LoadCommand(appBearer core.AppBearer, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer, apiError module.ApiError, appModule module.App) error {
	return run(Registry(), appBearer, pluginBearer, dbBearer, apiError, appModule, func(p module.Plugin, ctx module.PluginContext) error {
		return p.LoadCommand(ctx)
	})
}

func run(registry []module.Plugin, appBearer core.AppBearer, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer, apiError module.ApiError, appModule module.App, invoke func(module.Plugin, module.PluginContext) error) error {
	active, err := activePlugins(registry, appBearer, pluginBearer)
	if err != nil {
		return err
	}

	ordered, err := topoSort(active)
	if err != nil {
		return err
	}

	for _, p := range ordered {
		version, verr := pluginBearer.PluginVersion(p.Name())
		if verr != nil {
			return fmt.Errorf("plugin %q: %w", p.Name(), verr)
		}

		ctx := buildContext(p, version, pluginBearer, dbBearer, appBearer, apiError, appModule)

		if err := invoke(p, ctx); err != nil {
			return fmt.Errorf("plugin %q: %w", p.Name(), err)
		}
	}

	return nil
}

// activePlugins applies the AND-gate. A plugin activated in app.yml but missing
// its config file is an error — that's an explicit intent the operator can't
// fulfill. The inverse (config file present, plugin not activated) is silent
// since operators may stage config ahead of rollout.
func activePlugins(registry []module.Plugin, appBearer core.AppBearer, pluginBearer core.PluginBearer) ([]module.Plugin, error) {
	var active []module.Plugin
	appConfig := appBearer.Config()

	for _, p := range registry {
		if !appConfig.PluginExists(p.Name()) {
			continue
		}
		if !pluginBearer.ConfigExists(p.Name()) {
			return nil, fmt.Errorf("plugin %q is listed in app.yml but has no config/plugin/%s.yml file", p.Name(), p.Name())
		}
		active = append(active, p)
	}

	return active, nil
}

// topoSort returns `plugins` in Kahn-algorithm topological order. Ties are
// broken alphabetically for reproducibility. DependsOn entries that are not
// in `plugins` (i.e. not active) are silently skipped: soft dependency
// semantics. Cycles are reported with the residual set.
func topoSort(plugins []module.Plugin) ([]module.Plugin, error) {
	index := make(map[string]module.Plugin, len(plugins))
	for _, p := range plugins {
		index[p.Name()] = p
	}

	indegree := make(map[string]int, len(plugins))
	for _, p := range plugins {
		if _, seen := indegree[p.Name()]; !seen {
			indegree[p.Name()] = 0
		}
		for _, dep := range p.DependsOn() {
			if _, ok := index[dep]; !ok {
				continue // soft dep: skip deps that aren't active
			}
			if dep == p.Name() {
				return nil, fmt.Errorf("plugin %q declares itself as a dependency", p.Name())
			}
			indegree[p.Name()]++
		}
	}

	// Ready set, sorted for deterministic output.
	ready := make([]string, 0, len(plugins))
	for name, d := range indegree {
		if d == 0 {
			ready = append(ready, name)
		}
	}
	sort.Strings(ready)

	ordered := make([]module.Plugin, 0, len(plugins))
	for len(ready) > 0 {
		next := ready[0]
		ready = ready[1:]
		ordered = append(ordered, index[next])

		// Decrement indegree of everything that depends on `next`. Because
		// DependsOn is directional (p depends on dep), we iterate every
		// plugin and look for `next` in its DependsOn list.
		var newlyReady []string
		for _, p := range plugins {
			for _, dep := range p.DependsOn() {
				if dep != next {
					continue
				}
				if _, ok := index[dep]; !ok {
					continue
				}
				indegree[p.Name()]--
				if indegree[p.Name()] == 0 {
					newlyReady = append(newlyReady, p.Name())
				}
			}
		}
		sort.Strings(newlyReady)
		ready = append(ready, newlyReady...)
		sort.Strings(ready)
	}

	if len(ordered) != len(plugins) {
		var residual []string
		for name, d := range indegree {
			if d > 0 {
				residual = append(residual, name)
			}
		}
		sort.Strings(residual)
		return nil, fmt.Errorf("cycle detected among plugins: %v", residual)
	}

	return ordered, nil
}

func buildContext(p module.Plugin, version string, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer, appBearer core.AppBearer, apiError module.ApiError, appModule module.App) module.PluginContext {
	name := p.Name()
	return module.PluginContext{
		Version:   version,
		AppBearer: appBearer,
		AppModule: appModule,
		ApiError:  apiError,
		Logger:    log.With().Array("tags", zerolog.Arr().Str("altair").Str("plugin").Str(name)).Logger(),
		DecodeConfig: func(target interface{}) error {
			return pluginBearer.DecodeConfig(name, target)
		},
		Database: func(instance string) (db.DB, core.DatabaseConfig, error) {
			return dbBearer.Database(instance)
		},
	}
}
