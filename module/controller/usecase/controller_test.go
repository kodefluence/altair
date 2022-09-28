package usecase_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/altair/module/controller/usecase"
	metricDummyUsecase "github.com/kodefluence/altair/plugin/metric/module/dummy/usecase"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type ControllerSuiteTest struct {
	mockCtrl     *gomock.Controller
	controller   *usecase.Controller
	metric       usecase.Metric
	httpInjector usecase.HttpInjector
	apierror     module.ApiError
	apiengine    *gin.Engine

	suite.Suite
}

func (suite *ControllerSuiteTest) SetupTest() {

	suite.mockCtrl = gomock.NewController(suite.T())
	suite.metric = metricDummyUsecase.NewDummy()
	suite.apierror = apierror.Provide()

	gin.SetMode(gin.ReleaseMode)
	suite.apiengine = gin.New()
	suite.httpInjector = suite.apiengine.Handle

	suite.controller = usecase.NewController(suite.httpInjector, suite.apierror, suite.metric, &cobra.Command{})
}

func (suite *ControllerSuiteTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *ControllerSuiteTest) Subtest(testcase string, subtest func()) {
	suite.SetupTest()
	suite.Run(testcase, subtest)
	suite.TearDownTest()
}
