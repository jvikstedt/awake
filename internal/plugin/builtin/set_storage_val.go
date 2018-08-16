package builtin

import (
	"github.com/jvikstedt/awake"
)

type SetStorageVal struct{}

func (SetStorageVal) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_set_storage_val",
		DisplayName: "Set Storage Value",
	}
}

func (SetStorageVal) Perform(scope awake.Scope) error {
	vars := scope.Variables()

	for key, val := range vars {
		scope.SetStorageVariable(key, val)
	}

	return nil
}
