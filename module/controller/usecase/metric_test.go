package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MetricSuiteTest struct {
	*ControllerSuiteTest
}

type fakeMetric struct{}

func (*fakeMetric) InjectCounter(metricName string, labels ...string)     {}
func (*fakeMetric) InjectHistogram(metricName string, labels ...string)   {}
func (*fakeMetric) Inc(metricName string, labels map[string]string) error { return nil }
func (*fakeMetric) Observe(metricName string, value float64, labels map[string]string) error {
	return nil
}

func TestMetric(t *testing.T) {
	suite.Run(t, &MetricSuiteTest{
		&ControllerSuiteTest{},
	})
}

func (suite *HttpSuiteTest) TestListMetric() {
	suite.controller.InjectMetric(&fakeMetric{}, &fakeMetric{}, &fakeMetric{}, &fakeMetric{})
	suite.Assert().Equal(4, len(suite.controller.ListMetric()))
}

func (suite *HttpSuiteTest) TestInjectMetric() {
	suite.controller.InjectMetric(&fakeMetric{}, &fakeMetric{}, &fakeMetric{}, &fakeMetric{})
}
