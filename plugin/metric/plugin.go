package metric

import (
	_ "embed"
	"fmt"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/metric/entity"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus"
)

//go:embed config.sample.yml
var sampleConfig []byte

// Plugin implements module.Plugin for the metric plugin. The metric plugin
// does not own any database schema, so Migrations returns nil and
// RequiresDatabase is implicitly false.
type Plugin struct{}

// Name implements module.Plugin.
func (*Plugin) Name() string { return "metric" }

// DependsOn implements module.Plugin. Metric has no hard dependencies.
func (*Plugin) DependsOn() []string { return nil }

// Migrations implements module.Plugin. Metric owns no schema.
func (*Plugin) Migrations(ctx module.PluginContext) []module.MigrationSet { return nil }

// SampleConfig implements module.Plugin.
func (*Plugin) SampleConfig() []byte { return sampleConfig }

// Load implements module.Plugin. It dispatches on PluginContext.Version and
// the parsed `config.provider` field.
func (*Plugin) Load(ctx module.PluginContext) error {
	switch ctx.Version {
	case "1.0":
		return loadV1_0(ctx)
	default:
		return fmt.Errorf("undefined template version: %s for metric plugin", ctx.Version)
	}
}

// LoadCommand implements module.Plugin. Metric exposes no CLI subcommands today.
func (*Plugin) LoadCommand(ctx module.PluginContext) error { return nil }

func loadV1_0(ctx module.PluginContext) error {
	var metricPlugin entity.MetricPlugin
	if err := ctx.DecodeConfig(&metricPlugin); err != nil {
		return err
	}

	switch metricPlugin.Config.Provider {
	case "prometheus":
		prometheus.Load(ctx.AppModule)
		return nil
	default:
		return fmt.Errorf("metric plugin provider `%s` is not supported", metricPlugin.Config.Provider)
	}
}
