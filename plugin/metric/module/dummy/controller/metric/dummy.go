package metric

// DummyMetric
type DummyMetric struct{}

// New DummyMetric for metric placeholder
func NewDummy() *DummyMetric {
	return &DummyMetric{}
}

func (p *DummyMetric) InjectCounter(metricName string, labels ...string)     {}
func (p *DummyMetric) InjectHistogram(metricName string, labels ...string)   {}
func (p *DummyMetric) Inc(metricName string, labels map[string]string) error { return nil }
func (p *DummyMetric) Observe(metricName string, value float64, labels map[string]string) error {
	return nil
}
