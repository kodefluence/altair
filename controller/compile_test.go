package controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/controller"
	"github.com/kodefluence/altair/core"
	metricDummyUsecase "github.com/kodefluence/altair/plugin/metric/module/dummy/usecase"
	"github.com/kodefluence/altair/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	t.Run("Run gracefully", func(t *testing.T) {
		t.Run("Perform GET request", func(t *testing.T) {
			gracefullController := NewFakeController("/graceful", "GET", func(c *gin.Context) {
				c.String(http.StatusOK, "%s", "OK")
			})

			controller.Compile(engine, metricDummyUsecase.NewDummy(), gracefullController)
			w := testhelper.PerformRequest(engine, gracefullController.Method(), gracefullController.Path(), nil)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "OK", w.Body.String())
		})

		t.Run("Perform POST request", func(t *testing.T) {
			t.Run("Get body is not error", func(t *testing.T) {
				gracefullPostController := NewFakeController("/graceful_post", "POST", func(c *gin.Context) {
					c.String(http.StatusOK, "%s", "OK")
				})

				controller.Compile(engine, metricDummyUsecase.NewDummy(), gracefullPostController)
				w := testhelper.PerformRequest(engine, gracefullPostController.Method(), gracefullPostController.Path(), strings.NewReader(`{"object":"testing"}`))

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "OK", w.Body.String())
			})

			t.Run("Get body error", func(t *testing.T) {
				gracefullPostBodyErrorController := NewFakeController("/graceful_post_body_error", "POST", func(c *gin.Context) {
					c.String(http.StatusOK, "%s", "OK")
				})

				controller.Compile(engine, metricDummyUsecase.NewDummy(), gracefullPostBodyErrorController)
				w := testhelper.PerformRequest(engine, gracefullPostBodyErrorController.Method(), gracefullPostBodyErrorController.Path(), errRequestReader{})

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "OK", w.Body.String())
			})
		})
	})

	t.Run("Run not gracefully", func(t *testing.T) {
		t.Run("Controller return status >= 400", func(t *testing.T) {
			notGracefullController := NewFakeController("/not_graceful", "GET", func(c *gin.Context) {
				c.String(http.StatusInternalServerError, "%s", "Are you kidding me? The server is just crash!")
			})

			controller.Compile(engine, metricDummyUsecase.NewDummy(), notGracefullController)
			w := testhelper.PerformRequest(engine, notGracefullController.Method(), notGracefullController.Path(), nil)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Equal(t, "Are you kidding me? The server is just crash!", w.Body.String())
		})

		t.Run("Controller panic and compiler recover", func(t *testing.T) {
			t.Run("Panic is a string", func(t *testing.T) {
				panicStringController := NewFakeController("/panic_string", "GET", func(c *gin.Context) {
					panic("Panic with string")
				})

				controller.Compile(engine, metricDummyUsecase.NewDummy(), panicStringController)
				w := testhelper.PerformRequest(engine, panicStringController.Method(), panicStringController.Path(), nil)

				var response responseExample
				err := json.Unmarshal([]byte(w.Body.String()), &response)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Nil(t, err, "error should be nil")
			})

			t.Run("Panic is an error", func(t *testing.T) {

				panicErrorController := NewFakeController("/panic_error", "GET", func(c *gin.Context) {
					panic(errors.New("Panic with an error"))
				})

				controller.Compile(engine, metricDummyUsecase.NewDummy(), panicErrorController)
				w := testhelper.PerformRequest(engine, panicErrorController.Method(), panicErrorController.Path(), nil)

				var response responseExample
				err := json.Unmarshal([]byte(w.Body.String()), &response)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Nil(t, err, "error should be nil")
			})

			t.Run("Panic is neither error or string", func(t *testing.T) {
				panicOtherController := NewFakeController("/panic_other", "GET", func(c *gin.Context) {
					panic(responseExample{})
				})

				controller.Compile(engine, metricDummyUsecase.NewDummy(), panicOtherController)
				w := testhelper.PerformRequest(engine, panicOtherController.Method(), panicOtherController.Path(), nil)

				var response responseExample
				err := json.Unmarshal([]byte(w.Body.String()), &response)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Nil(t, err, "error should be nil")
			})
		})
	})
}

type (
	fakeController struct {
		path        string
		mockHandler func(c *gin.Context)
		method      string
	}

	responseExample struct {
		Message string `json:"message"`
		Meta    gin.H  `json:"meta"`
	}

	errRequestReader struct{}
)

func (errRequestReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("unexpected error")
}

func (fc fakeController) Path() string {
	return fc.path
}

func (fc fakeController) Method() string {
	return fc.method
}

func (fc fakeController) Control(c *gin.Context) {
	fc.mockHandler(c)
}

func NewFakeController(path, method string, mockHandler func(c *gin.Context)) core.Controller {
	return fakeController{
		path:        path,
		mockHandler: mockHandler,
		method:      method,
	}
}
