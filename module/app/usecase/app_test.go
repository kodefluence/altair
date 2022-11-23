package usecase_test

import (
	"testing"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/app/usecase"
	"github.com/stretchr/testify/assert"
)

type fakeController struct{}

func (*fakeController) InjectMetric(http ...module.MetricController)               {}
func (*fakeController) InjectHTTP(http ...module.HttpController)                   {}
func (*fakeController) InjectCommand(command ...module.CommandController)          {}
func (*fakeController) InjectDownstream(downstream ...module.DownstreamController) {}
func (*fakeController) ListDownstream() []module.DownstreamController              { return nil }
func (*fakeController) ListMetric() []module.MetricController                      { return nil }

func TestApp(t *testing.T) {
	fakeCtrl := &fakeController{}
	app := usecase.NewApp(fakeCtrl)
	assert.Equal(t, fakeCtrl, app.Controller())
}
