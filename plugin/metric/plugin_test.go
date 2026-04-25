package metric_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/metric"
	"github.com/kodefluence/altair/plugin/metric/entity"
)

// Assumption: metric plugin's identifier is "metric".
func TestPluginName_IsMetric(t *testing.T) {
	assert.Equal(t, "metric", (&metric.Plugin{}).Name())
}

// Assumption: metric has no hard dependencies.
func TestPluginDependsOn_IsNil(t *testing.T) {
	assert.Nil(t, (&metric.Plugin{}).DependsOn())
}

// Assumption: SampleConfig returns the embedded sample with metric plugin
// markers.
func TestPluginSampleConfig_ContainsPluginAndVersion(t *testing.T) {
	got := string((&metric.Plugin{}).SampleConfig())
	assert.Contains(t, got, "plugin: metric")
	assert.Contains(t, got, `version: "1.0"`)
	assert.Contains(t, got, "prometheus")
}

// Assumption: metric owns no schema.
func TestPluginMigrations_AlwaysNil(t *testing.T) {
	assert.Nil(t, (&metric.Plugin{}).Migrations(module.PluginContext{}))
}

// Assumption: LoadCommand is always a no-op success.
func TestPluginLoadCommand_AlwaysNil(t *testing.T) {
	assert.Nil(t, (&metric.Plugin{}).LoadCommand(module.PluginContext{Version: "anything"}))
}

// Assumption: Load with unknown version returns a clear error including the
// version string and "metric".
func TestPluginLoad_RejectsUnknownVersion(t *testing.T) {
	err := (&metric.Plugin{}).Load(module.PluginContext{Version: "9.9"})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "9.9")
	assert.Contains(t, err.Error(), "metric")
}

// Assumption: Load("1.0") with nil DecodeConfig errors out, doesn't panic.
// Mirrors the oauth plugin's defensive guard.
func TestPluginLoad_V10WithMissingDecodeConfigErrors(t *testing.T) {
	err := (&metric.Plugin{}).Load(module.PluginContext{Version: "1.0"})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "DecodeConfig")
}

// Assumption: Load("1.0") with a DecodeConfig that fails propagates that error.
func TestPluginLoad_V10DecodeConfigErrorPropagated(t *testing.T) {
	ctx := module.PluginContext{
		Version: "1.0",
		DecodeConfig: func(_ interface{}) error {
			return errors.New("decode boom")
		},
	}
	err := (&metric.Plugin{}).Load(ctx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "decode boom")
}

// Assumption: Load("1.0") with provider="" or unknown returns an error
// naming the bad provider.
func TestPluginLoad_V10UnsupportedProviderErrors(t *testing.T) {
	ctx := module.PluginContext{
		Version: "1.0",
		DecodeConfig: func(target interface{}) error {
			cfg, ok := target.(*entity.MetricPlugin)
			if !ok {
				t.Fatalf("expected *entity.MetricPlugin, got %T", target)
			}
			cfg.Config.Provider = "nope"
			return nil
		},
	}
	err := (&metric.Plugin{}).Load(ctx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "nope")
	assert.Contains(t, err.Error(), "not supported")
}
