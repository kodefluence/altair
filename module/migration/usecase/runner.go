package usecase

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
)

// PluginMigrationStatus is the report shape for `altair plugin migrate:status`
// and the drift-detection path that runs on `altair run`.
type PluginMigrationStatus struct {
	Plugin           string
	DatabaseInstance string
	Driver           string
	VersionTable     string
	HasMigrations    bool
	CurrentVersion   uint
	CurrentDirty     bool
	TargetVersion    uint
}

// DriftReport names a plugin whose embedded migrations are ahead of the DB.
// Returned by Runner.Drift() for the boot-time warning in altair run.
type DriftReport struct {
	Plugin           string
	DatabaseInstance string
	CurrentVersion   uint
	TargetVersion    uint
	Dirty            bool
}

// ErrPluginNotInRegistry is returned when a named plugin is not compiled in.
var ErrPluginNotInRegistry = errors.New("plugin not in registry")

// Runner drives migrations for every plugin in the registry. It is safe to
// call from both `altair run` (auto-migrate path) and `altair plugin migrate:*`
// (one-shot CLI path). Runner intentionally never calls migrator.Close():
// the golang-migrate mysql driver's Close() closes the underlying *sql.DB,
// which is the shared API-server connection pool.
type Runner struct {
	registry     []module.Plugin
	pluginBearer core.PluginBearer
	dbBearer     core.DatabaseBearer
}

// NewRunner constructs a Runner. pluginBearer is required because Plugin.Migrations()
// is allowed to depend on the plugin's parsed config (e.g. to learn which
// database instance to migrate).
func NewRunner(registry []module.Plugin, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer) *Runner {
	return &Runner{
		registry:     registry,
		pluginBearer: pluginBearer,
		dbBearer:     dbBearer,
	}
}

// UpAll runs Up on every plugin with migrations, in registry order. Stops at
// the first failure without rolling back earlier success — schema rollbacks
// across plugin boundaries are the operator's call, not the framework's.
func (r *Runner) UpAll() error {
	for _, p := range r.registry {
		if !r.configured(p) {
			continue
		}
		if err := r.Up(p.Name()); err != nil {
			return err
		}
	}
	return nil
}

// Up applies all pending migrations for a single plugin.
func (r *Runner) Up(pluginName string) error {
	sets, err := r.pluginMigrationSets(pluginName)
	if err != nil {
		return err
	}
	for _, set := range sets {
		if err := r.runOne(pluginName, set, func(m *migrate.Migrate) error {
			err := m.Up()
			if errors.Is(err, migrate.ErrNoChange) {
				return nil
			}
			return err
		}); err != nil {
			return fmt.Errorf("plugin %q: %w", pluginName, err)
		}
	}
	return nil
}

// Down reverses the last `steps` migrations for a single plugin. steps must
// be positive.
func (r *Runner) Down(pluginName string, steps int) error {
	if steps <= 0 {
		steps = 1
	}
	sets, err := r.pluginMigrationSets(pluginName)
	if err != nil {
		return err
	}
	for _, set := range sets {
		if err := r.runOne(pluginName, set, func(m *migrate.Migrate) error {
			err := m.Steps(-steps)
			if errors.Is(err, migrate.ErrNoChange) {
				return nil
			}
			return err
		}); err != nil {
			return fmt.Errorf("plugin %q: %w", pluginName, err)
		}
	}
	return nil
}

// Force sets the schema_version for a plugin, used to clear dirty state after
// a failed manual fix. Destructive: does not run SQL.
func (r *Runner) Force(pluginName string, version int) error {
	sets, err := r.pluginMigrationSets(pluginName)
	if err != nil {
		return err
	}
	for _, set := range sets {
		if err := r.runOne(pluginName, set, func(m *migrate.Migrate) error {
			return m.Force(version)
		}); err != nil {
			return fmt.Errorf("plugin %q: %w", pluginName, err)
		}
	}
	return nil
}

// Status returns one PluginMigrationStatus per plugin that has migrations.
// Plugins without embedded migrations are reported with HasMigrations=false.
// Plugins not yet configured (no plugin bearer entry) are skipped entirely.
func (r *Runner) Status() ([]PluginMigrationStatus, error) {
	var out []PluginMigrationStatus
	for _, p := range r.registry {
		if !r.configured(p) {
			continue
		}
		sets, err := r.pluginMigrationSets(p.Name())
		if err != nil {
			return nil, err
		}
		if len(sets) == 0 {
			out = append(out, PluginMigrationStatus{Plugin: p.Name(), HasMigrations: false})
			continue
		}
		for _, set := range sets {
			status := PluginMigrationStatus{
				Plugin:           p.Name(),
				DatabaseInstance: set.DatabaseInstance,
				Driver:           set.Driver,
				VersionTable:     r.resolveVersionTable(p.Name(), set),
				HasMigrations:    true,
			}
			if err := r.runOne(p.Name(), set, func(m *migrate.Migrate) error {
				if v, dirty, verr := m.Version(); verr == nil {
					status.CurrentVersion = v
					status.CurrentDirty = dirty
				} else if !errors.Is(verr, migrate.ErrNilVersion) {
					return verr
				}
				status.TargetVersion = maxSourceVersion(set.FS, set.SourcePath)
				return nil
			}); err != nil {
				return nil, fmt.Errorf("plugin %q: %w", p.Name(), err)
			}
			out = append(out, status)
		}
	}
	return out, nil
}

// Drift returns one DriftReport per plugin whose embedded schema is ahead of
// the DB. Callers use this for the on-boot warning in `altair run`.
func (r *Runner) Drift() ([]DriftReport, error) {
	statuses, err := r.Status()
	if err != nil {
		return nil, err
	}
	var reports []DriftReport
	for _, s := range statuses {
		if !s.HasMigrations {
			continue
		}
		if s.CurrentDirty || s.CurrentVersion < s.TargetVersion {
			reports = append(reports, DriftReport{
				Plugin:           s.Plugin,
				DatabaseInstance: s.DatabaseInstance,
				CurrentVersion:   s.CurrentVersion,
				TargetVersion:    s.TargetVersion,
				Dirty:            s.CurrentDirty,
			})
		}
	}
	return reports, nil
}

func (r *Runner) configured(p module.Plugin) bool {
	return r.pluginBearer.ConfigExists(p.Name())
}

func (r *Runner) pluginMigrationSets(pluginName string) ([]module.MigrationSet, error) {
	p, ok := r.findPlugin(pluginName)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrPluginNotInRegistry, pluginName)
	}
	ctx := module.PluginContext{
		DecodeConfig: func(target interface{}) error {
			return r.pluginBearer.DecodeConfig(pluginName, target)
		},
	}
	return p.Migrations(ctx), nil
}

func (r *Runner) findPlugin(name string) (module.Plugin, bool) {
	for _, p := range r.registry {
		if p.Name() == name {
			return p, true
		}
	}
	return nil, false
}

// runOne opens a migrator scoped to set, invokes action, and returns the
// resulting error. Never calls migrator.Close() — see Runner's doc comment.
func (r *Runner) runOne(pluginName string, set module.MigrationSet, action func(*migrate.Migrate) error) error {
	if set.DatabaseInstance == "" {
		return errors.New("migration set has empty DatabaseInstance")
	}

	sqldb, dbConfig, err := r.dbBearer.Database(set.DatabaseInstance)
	if err != nil {
		return err
	}

	versionTable := r.resolveVersionTable(pluginName, set)

	var dbDriver database.Driver
	switch set.Driver {
	case "mysql":
		dbDriver, err = mysql.WithInstance(sqldb.Eject(), &mysql.Config{
			MigrationsTable: versionTable,
			DatabaseName:    dbConfig.DBDatabase(),
		})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("migration driver %q is not supported", set.Driver)
	}

	sourceDriver, err := iofs.New(set.FS, set.SourcePath)
	if err != nil {
		return err
	}
	defer sourceDriver.Close()

	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, set.Driver, dbDriver)
	if err != nil {
		return err
	}
	// migrator.Close() would close dbDriver, which closes the shared *sql.DB.
	// Callers manage sql.DB lifetime via db.CloseAll() on process exit.

	return action(migrator)
}

func (r *Runner) resolveVersionTable(pluginName string, set module.MigrationSet) string {
	if set.VersionTable != "" {
		return set.VersionTable
	}
	return fmt.Sprintf("%s_plugin_db_versions", pluginName)
}
