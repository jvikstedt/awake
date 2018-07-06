package builtin

import (
	"fmt"

	"github.com/jvikstedt/awake"
)

type LessThan struct{}

func (LessThan) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_less_than",
		DisplayName: "Less Than",
	}
}

func (LessThan) Perform(scope awake.Scope) error {
	actual, _ := scope.ValueAsFloat("actual")
	expected, _ := scope.ValueAsFloat("expected")

	if actual > expected {
		return fmt.Errorf("Expected to be less than %v but got %v", expected, actual)
	}

	return nil
}
