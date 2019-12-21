package controller_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/core"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	t.Run("Run gracefully", func(t *testing.T) {
		t.Run("Perform GET request", func(t *testing.T) {
			gracefullController := NewFakeController("/gracefull", "GET", func(c *gin.Context) {
				c.String(http.StatusOK, "%s", "OK")
			})

			controller.Compile(engine, gracefullController)
			w := performRequest(engine, gracefullController.Method(), gracefullController.Path(), nil)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "OK", w.Body.String())
		})

		t.Run("Perform POST request", func(t *testing.T) {
			t.Run("Get body is not error", func(t *testing.T) {
				gracefullPostController := NewFakeController("/gracefull_post", "POST", func(c *gin.Context) {
					c.String(http.StatusOK, "%s", "OK")
				})

				controller.Compile(engine, gracefullPostController)
				w := performRequest(engine, gracefullPostController.Method(), gracefullPostController.Path(), strings.NewReader(`{"object":"testing"}`))

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "OK", w.Body.String())
			})

			t.Run("Get body error", func(t *testing.T) {
				gracefullPostBodyErrorController := NewFakeController("/gracefull_post_body_error", "POST", func(c *gin.Context) {
					c.String(http.StatusOK, "%s", "OK")
				})

				controller.Compile(engine, gracefullPostBodyErrorController)
				w := performRequest(engine, gracefullPostBodyErrorController.Method(), gracefullPostBodyErrorController.Path(), errRequestReader{})

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "OK", w.Body.String())
			})
		})
	})

	t.Run("Run not gracefully", func(t *testing.T) {
		t.Run("Controller return status >= 400", func(t *testing.T) {
			notGracefullController := NewFakeController("/not_gracefull", "GET", func(c *gin.Context) {
				c.String(http.StatusInternalServerError, "%s", "Are you kidding me? The server is just crash!")
			})

			controller.Compile(engine, notGracefullController)
			w := performRequest(engine, notGracefullController.Method(), notGracefullController.Path(), nil)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Equal(t, "Are you kidding me? The server is just crash!", w.Body.String())
		})

		t.Run("Controller panic and compiler recover", func(t *testing.T) {
			t.Run("Panic is a string", func(t *testing.T) {
				expectedBody := responseExample{
					Meta: gin.H{"http_status": http.StatusInternalServerError},
				}

				panicStringController := NewFakeController("/panic_string", "GET", func(c *gin.Context) {
					panic("Panic with string")
				})

				controller.Compile(engine, panicStringController)
				w := performRequest(engine, panicStringController.Method(), panicStringController.Path(), nil)

				var response responseExample
				err := json.Unmarshal([]byte(w.Body.String()), &response)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.EqualValues(t, expectedBody.Meta["http_status"], response.Meta["http_status"])
				assert.Nil(t, err, "error should be nil")
			})

			t.Run("Panic is an error", func(t *testing.T) {
				expectedBody := responseExample{
					Meta: gin.H{"http_status": http.StatusInternalServerError},
				}

				panicErrorController := NewFakeController("/panic_error", "GET", func(c *gin.Context) {
					panic(errors.New("Panic with an error"))
				})

				controller.Compile(engine, panicErrorController)
				w := performRequest(engine, panicErrorController.Method(), panicErrorController.Path(), nil)

				var response responseExample
				err := json.Unmarshal([]byte(w.Body.String()), &response)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.EqualValues(t, expectedBody.Meta["http_status"], response.Meta["http_status"])
				assert.Nil(t, err, "error should be nil")
			})

			t.Run("Panic is neither error or string", func(t *testing.T) {
				expectedBody := responseExample{
					Meta: gin.H{"http_status": http.StatusInternalServerError},
				}

				panicOtherController := NewFakeController("/panic_other", "GET", func(c *gin.Context) {
					panic(responseExample{})
				})

				controller.Compile(engine, panicOtherController)
				w := performRequest(engine, panicOtherController.Method(), panicOtherController.Path(), nil)

				var response responseExample
				err := json.Unmarshal([]byte(w.Body.String()), &response)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.EqualValues(t, expectedBody.Meta["http_status"], response.Meta["http_status"])
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

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
