package plugin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/altair/module"
)

// TestRegistry_DependsOnReferencesAreResolvable is a compile-time-ish audit:
// every name in any plugin's DependsOn() list must appear in Registry(),
// otherwise the topo sort will silently treat it as a soft-missing dep and
// the intended ordering constraint quietly decays.
func TestRegistry_DependsOnReferencesAreResolvable(t *testing.T) {
	names := map[string]bool{}
	for _, p := range Registry() {
		names[p.Name()] = true
	}

	for _, p := range Registry() {
		for _, dep := range p.DependsOn() {
			assert.True(t, names[dep], "plugin %q DependsOn unknown plugin %q", p.Name(), dep)
		}
	}
}

// TestRegistry_NamesAreUnique makes the registry's uniqueness invariant
// explicit. Two plugins with the same Name() would collide on every
// pluginBearer lookup.
func TestRegistry_NamesAreUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range Registry() {
		assert.False(t, seen[p.Name()], "duplicate plugin name in Registry(): %q", p.Name())
		seen[p.Name()] = true
	}
}

// TestRegistry_MigrationTablesAreUnique enforces that no two plugins declare
// the same (DatabaseInstance, VersionTable) tuple. golang-migrate interleaves
// unrelated schemas under a shared migrations table as if they were a single
// linear history, corrupting state. This test uses a stub DecodeConfig so
// Migrations(ctx) resolves every plugin's database instance to a predictable
// string; real operators will configure different instances, but the
// invariant is that the *declared* (instance, table) pair must be unique per
// plugin regardless of which instance gets plugged in at runtime.
func TestRegistry_MigrationTablesAreUnique(t *testing.T) {
	seen := map[string]string{} // key = "<instance>|<table>", value = owning plugin
	for _, p := range Registry() {
		ctx := module.PluginContext{
			DecodeConfig: func(target interface{}) error {
				// Many plugins ignore decode errors and return nil migrations;
				// that's fine for this audit.
				return nil
			},
		}
		for _, set := range p.Migrations(ctx) {
			// VersionTable defaults to "<name>_plugin_db_versions" if empty;
			// mirror the runner's resolution so we compare effective values.
			vt := set.VersionTable
			if vt == "" {
				vt = fmt.Sprintf("%s_plugin_db_versions", p.Name())
			}
			key := fmt.Sprintf("%s|%s", set.DatabaseInstance, vt)
			if owner, clash := seen[key]; clash {
				t.Fatalf("plugins %q and %q both declare migration table %q on instance %q", owner, p.Name(), vt, set.DatabaseInstance)
			}
			seen[key] = p.Name()
		}
	}
}
