package builtin

import (
	"fmt"

	"github.com/jvikstedt/awake"
)

type MoreThan struct{}

func (MoreThan) Tag() string {
	return "builtin_more_than"
}

func (MoreThan) Perform(scope awake.Scope) error {
	actual, _ := scope.ValueAsFloat("actual")
	expected, _ := scope.ValueAsFloat("expected")

	if actual < expected {
		return fmt.Errorf("Expected to be more than %v but got %v", expected, actual)
	}

	return nil
}
