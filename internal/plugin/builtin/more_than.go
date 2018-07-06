package builtin

import (
	"fmt"

	"github.com/jvikstedt/awake"
)

type MoreThan struct{}

func (MoreThan) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_more_than",
		DisplayName: "More Than",
	}
}

func (MoreThan) Perform(scope awake.Scope) error {
	actual, _ := scope.ValueAsFloat("actual")
	expected, _ := scope.ValueAsFloat("expected")

	if actual < expected {
		return fmt.Errorf("Expected to be more than %v but got %v", expected, actual)
	}

	return nil
}
