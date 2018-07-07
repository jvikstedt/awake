package builtin

import (
	"github.com/jvikstedt/awake"
)

type Errors struct{}

func (Errors) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_errors",
		DisplayName: "Error count",
	}
}

func (Errors) Perform(scope awake.Scope) error {
	errors := scope.Errors()

	scope.SetReturnVariable("count", awake.Variable{Type: awake.TypeInt, Val: len(errors)})
	scope.SetReturnVariable("hasErrors", awake.Variable{Type: awake.TypeBool, Val: len(errors) > 0})

	errStr := ""
	for _, e := range errors {
		errStr += e.Error() + "\n"
	}

	scope.SetReturnVariable("errors", awake.Variable{Type: awake.TypeString, Val: errStr})

	return nil
}
