package usecase_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/entity"
	"github.com/stretchr/testify/suite"
)

type DownstreamSuiteTest struct {
	*ControllerSuiteTest
}

type fakeDownstream struct{}

func (*fakeDownstream) Name() string {
	return "fake-downstream"
}

func (*fakeDownstream) Intervene(c *gin.Context, proxyReq *http.Request, r entity.RouterPath) error {
	return nil
}

func TestDownstream(t *testing.T) {
	suite.Run(t, &DownstreamSuiteTest{
		&ControllerSuiteTest{},
	})
}

func (suite *HttpSuiteTest) TestListDownstream() {
	suite.controller.InjectDownstream(&fakeDownstream{}, &fakeDownstream{}, &fakeDownstream{}, &fakeDownstream{})
	suite.Assert().Equal(4, len(suite.controller.ListDownstream()))
}

func (suite *HttpSuiteTest) TestInjectDownstream() {
	suite.controller.InjectDownstream(&fakeDownstream{}, &fakeDownstream{}, &fakeDownstream{}, &fakeDownstream{})
}
