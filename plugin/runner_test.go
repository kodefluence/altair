package plugin

import (
	"strings"
	"testing"

	"github.com/kodefluence/altair/module"
	"github.com/stretchr/testify/assert"
)

type stubPlugin struct {
	name string
	deps []string
}

func (s *stubPlugin) Name() string                                          { return s.name }
func (s *stubPlugin) DependsOn() []string                                   { return s.deps }
func (s *stubPlugin) Migrations(module.PluginContext) []module.MigrationSet { return nil }
func (s *stubPlugin) Load(module.PluginContext) error                       { return nil }
func (s *stubPlugin) LoadCommand(module.PluginContext) error                { return nil }
func (s *stubPlugin) SampleConfig() []byte                                  { return nil }

func names(plugins []module.Plugin) []string {
	out := make([]string, len(plugins))
	for i, p := range plugins {
		out[i] = p.Name()
	}
	return out
}

func TestTopoSort_AlphabeticalTieBreak(t *testing.T) {
	// No dependencies; alphabetical ordering keeps `make test` diffs stable.
	plugins := []module.Plugin{
		&stubPlugin{name: "zeta"},
		&stubPlugin{name: "alpha"},
		&stubPlugin{name: "mu"},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, err)
	assert.Equal(t, []string{"alpha", "mu", "zeta"}, names(ordered))
}

func TestTopoSort_RespectsDependsOn(t *testing.T) {
	plugins := []module.Plugin{
		&stubPlugin{name: "oauth", deps: []string{"metric"}},
		&stubPlugin{name: "metric"},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, err)
	assert.Equal(t, []string{"metric", "oauth"}, names(ordered))
}

func TestTopoSort_SkipsSoftMissingDependency(t *testing.T) {
	// oauth depends on metric, but metric is not in the active set;
	// soft dep semantics say we load oauth anyway without error.
	plugins := []module.Plugin{
		&stubPlugin{name: "oauth", deps: []string{"metric"}},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, err)
	assert.Equal(t, []string{"oauth"}, names(ordered))
}

func TestTopoSort_RejectsCycles(t *testing.T) {
	plugins := []module.Plugin{
		&stubPlugin{name: "a", deps: []string{"b"}},
		&stubPlugin{name: "b", deps: []string{"a"}},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, ordered)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "cycle"), "expected cycle error, got: %v", err)
	assert.True(t, strings.Contains(err.Error(), "a"), "expected residual 'a' in error")
	assert.True(t, strings.Contains(err.Error(), "b"), "expected residual 'b' in error")
}

func TestTopoSort_RejectsSelfDependency(t *testing.T) {
	plugins := []module.Plugin{
		&stubPlugin{name: "loner", deps: []string{"loner"}},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, ordered)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "loner"), "expected plugin name in self-dep error")
}
