package testhelper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/kodefluence/altair/module/apierror"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

type MockErrorIoReader struct {
}

func (m MockErrorIoReader) Read(x []byte) (int, error) {
	return 0, errors.New("read error")
}

func PerformRequest(r http.Handler, method, path string, body io.Reader, reqModifiers ...func(req *http.Request)) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	for _, f := range reqModifiers {
		f(req)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func GenerateTempTestFiles(configPath, content, fileName string, mode os.FileMode) {
	err := os.Mkdir(configPath, os.ModePerm)
	if err != nil {
		if pathError, ok := err.(*os.PathError); ok && pathError.Err.Error() != "file exists" {
			panic(err)
		}
	}

	f, err := os.OpenFile(fmt.Sprintf("%s%s", configPath, fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func RemoveTempTestFiles(configPath string) {
	err := os.RemoveAll(configPath)
	if err != nil {
		panic(err)
	}
}

func ErrInternalServer() jsonapi.Errors {
	return jsonapi.BuildResponse(apierror.Provide().InternalServerError(kontext.Fabricate())).Errors
}
