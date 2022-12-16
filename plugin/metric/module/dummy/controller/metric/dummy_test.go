package metric_test

import (
	"testing"

	"github.com/kodefluence/altair/plugin/metric/module/dummy/controller/metric"
	"github.com/stretchr/testify/assert"
)

func TestPrometheus(t *testing.T) {
	dummyMetric := metric.NewDummy()

	t.Run("InjectCounter", func(t *testing.T) {
		dummyMetric.InjectCounter("testing_metrics", "blablabla")
	})

	t.Run("InjectHistogram", func(t *testing.T) {
		dummyMetric.InjectHistogram("testing_metrics", "blablabla")
	})

	t.Run("Inc", func(t *testing.T) {
		err := dummyMetric.Inc("testing_metrics", make(map[string]string))
		assert.Nil(t, err)
	})

	t.Run("Observer", func(t *testing.T) {
		err := dummyMetric.Observe("testing_metric", 0, make(map[string]string))
		assert.Nil(t, err)
	})
}
