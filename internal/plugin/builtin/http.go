package builtin

import (
	"io/ioutil"
	"net/http"

	"github.com/jvikstedt/awake"
)

type HTTP struct{}

func (HTTP) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_http",
		DisplayName: "HTTP",
	}
}

func (HTTP) Perform(scope awake.Scope) error {
	url, _ := scope.ValueAsString("url")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	scope.SetReturnVariable("code", awake.Variable{
		Type: awake.TypeInt,
		Val:  resp.StatusCode,
	})

	scope.SetReturnVariable("body", awake.Variable{
		Type: awake.TypeArrayBytes,
		Val:  bodyBytes,
	})

	return nil
}
