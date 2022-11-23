package metric_test

import (
	"testing"

	"github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/metric"
	"github.com/stretchr/testify/assert"
)

func TestPrometheus(t *testing.T) {
	promMetric := metric.NewPrometheus()

	t.Run("InjectCounter", func(t *testing.T) {
		t.Run("Metric is not exists", func(t *testing.T) {
			promMetric.InjectCounter("some_metric")
		})

		t.Run("Metric already exists", func(t *testing.T) {
			promMetric.InjectCounter("some_metric")
		})
	})

	t.Run("InjectHistogram", func(t *testing.T) {
		t.Run("Metric is not exists", func(t *testing.T) {
			promMetric.InjectHistogram("some_metric_histogram")
		})

		t.Run("Metric already exists", func(t *testing.T) {
			promMetric.InjectHistogram("some_metric_histogram")
		})
	})

	t.Run("Inc", func(t *testing.T) {
		t.Run("Run gracefully", func(t *testing.T) {
			t.Run("Return nil", func(t *testing.T) {
				assert.Nil(t, promMetric.Inc("some_metric", nil))
			})
		})

		t.Run("Metric is not exists", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				assert.NotNil(t, promMetric.Inc("some_metric_that_not_exists", nil))
			})
		})

		t.Run("Get metric with labels", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				promMetric.InjectCounter("some_metric_with_labels", "label_a", "label_b")
				assert.NotNil(t, promMetric.Inc("some_metric_with_labels", nil))
			})
		})
	})

	t.Run("Observer", func(t *testing.T) {
		t.Run("Run gracefully", func(t *testing.T) {
			t.Run("Return nil", func(t *testing.T) {
				assert.Nil(t, promMetric.Observe("some_metric_histogram", 0, nil))
			})
		})

		t.Run("Metric is not exists", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				assert.NotNil(t, promMetric.Observe("some_metric_histogram_that_not_exists", 0, nil))
			})
		})

		t.Run("Get metric with labels", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				promMetric.InjectHistogram("some_metric_histogram_with_labels", "label_a", "label_b")
				assert.NotNil(t, promMetric.Observe("some_metric_histogram_with_labels", 0, nil))
			})
		})
	})
}
