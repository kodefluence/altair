package testhelper

// DummyMetric is a no-op module.MetricController used in tests that need a
// metric provider without asserting on the metric values. It lived in
// plugin/metric/module/dummy/ until the dummy fallback was inlined into the
// framework (no metric controller registered means no metric calls).
type DummyMetric struct{}

// NewDummyMetric returns a no-op metric controller.
func NewDummyMetric() *DummyMetric {
	return &DummyMetric{}
}

func (*DummyMetric) InjectCounter(metricName string, labels ...string)     {}
func (*DummyMetric) InjectHistogram(metricName string, labels ...string)   {}
func (*DummyMetric) Inc(metricName string, labels map[string]string) error { return nil }
func (*DummyMetric) Observe(metricName string, value float64, labels map[string]string) error {
	return nil
}
