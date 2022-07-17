package cfg

import (
	"bytes"
	"os"
	"text/template"

	"github.com/google/uuid"
)

func compileTemplate(b []byte) ([]byte, error) {
	tpl, err := template.New(uuid.New().String()).Funcs(template.FuncMap{
		"env": envFallback,
	}).Parse(string(b))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, nil)
	return buf.Bytes(), err
}

func envFallback(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		return `""`
	}

	return env
}
