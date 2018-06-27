package builtin

import (
	"github.com/jvikstedt/awake"
	"github.com/tidwall/gjson"
)

type JSON struct{}

func (j JSON) Tag() string {
	return "builtin_json"
}

func (j JSON) Perform(scope awake.Scope) error {
	json, _ := scope.ValueAsBytes("json")
	path, _ := scope.ValueAsString("path")
	returnName, _ := scope.ValueAsString("returnName")

	result := gjson.GetBytes(json, path)

	scope.SetReturnVariable(returnName, awake.Variable{Type: "unknown", Val: result.Value()})

	return nil
}
