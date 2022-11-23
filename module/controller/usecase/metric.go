package usecase

import "github.com/kodefluence/altair/module"

func (c *Controller) InjectMetric(metric ...module.MetricController) {
	c.metricController = append(c.metricController, metric...)
}

func (c *Controller) ListMetric() []module.MetricController {
	return c.metricController
}
