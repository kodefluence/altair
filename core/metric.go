package core

type Metric interface {
	InjectCounter(metricName string, labels ...string)
	InjectHistogram(metricName string, labels ...string)
	Inc(metricName string, labels map[string]string) error
	Observe(metricName string, value float64, labels map[string]string) error
}
