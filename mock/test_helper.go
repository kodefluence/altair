package mock

import (
	"io"
	"net/http"
	"net/http/httptest"

	entity "github.com/codefluence-x/altair/entity"
)

type ErrorResponse struct {
	Errors []entity.ErrorObject `json:"errors"`
}

func PerformRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
