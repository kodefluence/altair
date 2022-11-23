package usecase_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/plugin/metric/module/dummy/controller/metric"
	"github.com/kodefluence/altair/testhelper"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/suite"
)

type HttpSuiteTest struct {
	*ControllerSuiteTest
}

type fakeHttpController struct {
	path        string
	method      string
	mockHandler func(ktx kontext.Context, c *gin.Context)
}

func (fc fakeHttpController) Path() string                                { return fc.path }
func (fc fakeHttpController) Method() string                              { return fc.method }
func (fc fakeHttpController) Control(ktx kontext.Context, c *gin.Context) { fc.mockHandler(ktx, c) }

type errRequestReader struct{}

func (errRequestReader) Read(p []byte) (n int, err error) { return 0, errors.New("unexpected error") }

func TestHttp(t *testing.T) {
	suite.Run(t, &HttpSuiteTest{
		&ControllerSuiteTest{},
	})
}

func (suite *HttpSuiteTest) TestInjectHttp() {
	suite.Run("Positive cases", func() {
		suite.Subtest("Perform GET request", func() {
			fakecontroller := fakeHttpController{
				"/graceful", "GET", func(ktx kontext.Context, c *gin.Context) {
					c.String(http.StatusOK, "%s", "OK")
				}}

			suite.controller.InjectMetric(metric.NewDummy())
			suite.controller.InjectHTTP(fakecontroller)
			w := testhelper.PerformRequest(suite.apiengine, fakecontroller.Method(), fakecontroller.Path(), nil)

			suite.Assert().Equal(http.StatusOK, w.Code)
			suite.Assert().Equal("OK", w.Body.String())
		})

		suite.Subtest("Perform POST", func() {
			body := "this is beautifull body"

			fakecontroller := fakeHttpController{
				"/graceful", "POST", func(ktx kontext.Context, c *gin.Context) {
					reader, _ := c.Request.GetBody()
					requestBody, _ := io.ReadAll(reader)
					suite.Assert().Equal(body, string(requestBody))
					c.String(http.StatusOK, "%s", "OK")
				}}

			suite.controller.InjectHTTP(fakecontroller)
			w := testhelper.PerformRequest(suite.apiengine, fakecontroller.Method(), fakecontroller.Path(), strings.NewReader(body))

			suite.Assert().Equal(http.StatusOK, w.Code)
			suite.Assert().Equal("OK", w.Body.String())
		})

		suite.Subtest("Perform POST request with nil body", func() {
			fakecontroller := fakeHttpController{
				"/graceful", "POST", func(ktx kontext.Context, c *gin.Context) {
					c.String(http.StatusOK, "%s", "OK")
				}}

			suite.controller.InjectHTTP(fakecontroller)
			w := testhelper.PerformRequest(suite.apiengine, fakecontroller.Method(), fakecontroller.Path(), nil)

			suite.Assert().Equal(http.StatusOK, w.Code)
			suite.Assert().Equal("OK", w.Body.String())
		})

		suite.Subtest("Perform POST request with body error", func() {
			fakecontroller := fakeHttpController{
				"/graceful_with_body_error", "POST", func(ktx kontext.Context, c *gin.Context) {
					c.String(http.StatusOK, "%s", "OK")
				}}

			suite.controller.InjectHTTP(fakecontroller)
			w := testhelper.PerformRequest(suite.apiengine, fakecontroller.Method(), fakecontroller.Path(), errRequestReader{})

			suite.Assert().Equal(http.StatusOK, w.Code)
			suite.Assert().Equal("OK", w.Body.String())
		})
	})

	suite.Run("Negative cases", func() {
		suite.Subtest("http status >= 400", func() {
			fakecontroller := fakeHttpController{
				"/not_graceful", "GET", func(ktx kontext.Context, c *gin.Context) {
					c.String(http.StatusInternalServerError, "%s", "Are you kidding me? The server is just crash!")
				}}

			suite.controller.InjectHTTP(fakecontroller)
			w := testhelper.PerformRequest(suite.apiengine, fakecontroller.Method(), fakecontroller.Path(), errRequestReader{})

			suite.Assert().Equal(http.StatusInternalServerError, w.Code)
			suite.Assert().Equal("Are you kidding me? The server is just crash!", w.Body.String())
		})

		suite.Subtest("panic string and recover", func() {
			fakecontroller := fakeHttpController{
				"/panic_string", "GET", func(ktx kontext.Context, c *gin.Context) {
					panic("Panic with string")
				}}

			suite.controller.InjectHTTP(fakecontroller)
			w := testhelper.PerformRequest(suite.apiengine, fakecontroller.Method(), fakecontroller.Path(), errRequestReader{})

			suite.Assert().Equal(http.StatusInternalServerError, w.Code)
		})

		suite.Subtest("panic error and recover", func() {
			fakecontroller := fakeHttpController{
				"/panic_error", "GET", func(ktx kontext.Context, c *gin.Context) {
					panic(errors.New("Panic with an error"))
				}}

			suite.controller.InjectHTTP(fakecontroller)
			w := testhelper.PerformRequest(suite.apiengine, fakecontroller.Method(), fakecontroller.Path(), errRequestReader{})

			suite.Assert().Equal(http.StatusInternalServerError, w.Code)
		})

		suite.Subtest("panic neither error and string and recover", func() {
			fakecontroller := fakeHttpController{
				"/panic_neither_error_string", "GET", func(ktx kontext.Context, c *gin.Context) {
					panic(fakeHttpController{})
				}}

			suite.controller.InjectHTTP(fakecontroller)
			w := testhelper.PerformRequest(suite.apiengine, fakecontroller.Method(), fakecontroller.Path(), errRequestReader{})

			suite.Assert().Equal(http.StatusInternalServerError, w.Code)
		})
	})
}
