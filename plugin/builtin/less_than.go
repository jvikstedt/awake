package builtin

import (
	"fmt"

	"github.com/jvikstedt/awake"
)

type LessThan struct{}

func (l LessThan) Tag() string {
	return "builtin_less_than"
}

func (l LessThan) Perform(scope awake.Scope) error {
	actual, _ := scope.ValueAsFloat("actual")
	expected, _ := scope.ValueAsFloat("expected")

	if actual > expected {
		return fmt.Errorf("Expected to be less than %v but got %v", expected, actual)
	}

	return nil
}
