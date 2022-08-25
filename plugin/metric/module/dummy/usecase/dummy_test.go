package usecase_test

import (
	"testing"

	"github.com/kodefluence/altair/plugin/metric/module/dummy/usecase"
)

func TestPrometheus(t *testing.T) {
	dummyMetric := usecase.NewDummy()

	t.Run("InjectCounter", func(t *testing.T) {
		dummyMetric.InjectCounter("testing_metrics", "blablabla")
	})

	t.Run("InjectHistogram", func(t *testing.T) {
		dummyMetric.InjectHistogram("testing_metrics", "blablabla")
	})

	t.Run("Inc", func(t *testing.T) {
		dummyMetric.Inc("testing_metrics", make(map[string]string))
	})

	t.Run("Observer", func(t *testing.T) {
		dummyMetric.Observe("testing_metric", 0, make(map[string]string))
	})
}
