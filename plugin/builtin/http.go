package builtin

import (
	"net/http"

	"github.com/jvikstedt/awake"
)

type HTTP struct{}

func (h HTTP) Tag() string {
	return "builtin_http"
}

func (h HTTP) Perform(scope awake.Scope) error {
	url, _ := scope.ValueAsString("url")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	scope.SetReturnVariable("code", awake.Variable{
		Type: "integer",
		Val:  resp.StatusCode,
	})

	return nil
}
