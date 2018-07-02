package builtin

import (
	"fmt"

	"github.com/jvikstedt/awake"
)

type Equal struct{}

func (e Equal) Tag() string {
	return "builtin_equal"
}

func (e Equal) Perform(scope awake.Scope) error {
	actual, _ := scope.ValueAsRaw("actual")
	expected, _ := scope.ValueAsRaw("expected")

	if actual != expected {
		return fmt.Errorf("Expected to be %v but got %v", expected, actual)
	}

	return nil
}
