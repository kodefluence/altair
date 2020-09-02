package metric

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusMetric struct {
	counterMetrics    map[string]*prometheus.CounterVec
	counterMetricLock *sync.Mutex

	histogramMetrics    map[string]*prometheus.HistogramVec
	histogramMetricLock *sync.Mutex
}

func NewPrometheusMetric() *prometheusMetric {
	return &prometheusMetric{
		counterMetrics:    map[string]*prometheus.CounterVec{},
		counterMetricLock: &sync.Mutex{},

		histogramMetrics:    map[string]*prometheus.HistogramVec{},
		histogramMetricLock: &sync.Mutex{},
	}
}

func (p *prometheusMetric) InjectCounter(metricName string, labels ...string) {
	if _, ok := p.counterMetrics[metricName]; ok {
		return
	}

	p.counterMetricLock.Lock()
	counterMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metricName,
	}, labels)
	_ = prometheus.Register(counterMetric)
	p.counterMetrics[metricName] = counterMetric
	p.counterMetricLock.Unlock()
}

func (p *prometheusMetric) InjectHistogram(metricName string, labels ...string) {
	if _, ok := p.histogramMetrics[metricName]; ok {
		return
	}

	p.histogramMetricLock.Lock()
	histogramMetric := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: metricName,
	}, labels)
	_ = prometheus.Register(histogramMetric)
	p.histogramMetrics[metricName] = histogramMetric
	p.histogramMetricLock.Unlock()
}

func (p *prometheusMetric) Inc(metricName string, labels map[string]string) error {
	counterMetric, ok := p.counterMetrics[metricName]
	if !ok {
		return fmt.Errorf("Metric `%s` is not exists", metricName)
	}

	counter, err := counterMetric.GetMetricWith(labels)
	if err != nil {
		return err
	}

	counter.Inc()

	return nil
}

func (p *prometheusMetric) Observe(metricName string, value float64, labels map[string]string) error {
	histogramMetric, ok := p.histogramMetrics[metricName]
	if !ok {
		return fmt.Errorf("Metric `%s` is not exists", metricName)
	}

	observer, err := histogramMetric.GetMetricWith(labels)
	if err != nil {
		return err
	}

	observer.Observe(value)

	return nil
}
