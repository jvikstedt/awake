package builtin

import (
	"encoding/json"

	"github.com/jvikstedt/awake"
	"github.com/tidwall/gjson"
)

type JSON struct{}

func (JSON) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_json",
		DisplayName: "JSON",
	}
}

func (JSON) Perform(scope awake.Scope) error {
	data, _ := scope.ValueAsBytes("json")
	path, _ := scope.ValueAsString("path")
	returnName, _ := scope.ValueAsString("returnName")

	if path != "" {
		result := gjson.GetBytes(data, path)
		scope.SetReturnVariable(returnName, awake.MakeVariable(result.Value()))
		return nil
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	scope.SetReturnVariable(returnName, awake.MakeVariable(result))

	return nil
}
