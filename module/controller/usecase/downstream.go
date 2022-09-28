package usecase

import "github.com/kodefluence/altair/module"

func (c *Controller) ListDownstream() []module.DownstreamController {
	return c.downstreamController
}

func (c *Controller) InjectDownstream(downstream ...module.DownstreamController) {
	c.downstreamController = append(c.downstreamController, downstream...)
}
