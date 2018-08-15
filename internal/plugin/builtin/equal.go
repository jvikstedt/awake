package builtin

import (
	"fmt"

	"github.com/jvikstedt/awake"
)

type Equal struct{}

func (Equal) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_equal",
		DisplayName: "Equal",
	}
}

func (Equal) Perform(scope awake.Scope) error {
	actual, _ := scope.ValueAsRaw("actual")
	expected, _ := scope.ValueAsRaw("expected")

	if actual != expected {
		return fmt.Errorf("Expected to be %v %T but got %v %T", expected, expected, actual, actual)
	}

	return nil
}
