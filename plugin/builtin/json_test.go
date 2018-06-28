package builtin_test

import (
	"testing"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/plugin/builtin"
)

func TestTag(t *testing.T) {
	json := builtin.JSON{}
	if json.Tag() != "builtin_json" {
		t.Errorf("Expected tag to be builtin_json but got %s", json.Tag())
	}
}

func TestPerform(t *testing.T) {
}

type ScopeMock struct {
	Receives struct {
		ValueAsRaw struct {
			Name string
		}
		ValueAsString struct {
			Name string
		}
		ValueAsBytes struct {
			Name string
		}
		ValueAsInt struct {
			Name string
		}
		ValueAsFloat struct {
			Name string
		}
		ValueAsBool struct {
			Name string
		}
		SetReturnVariable struct {
			Name string
			awake.Variable
		}
	}
	Returns struct {
		ValueAsRaw struct {
			Val interface{}
			OK  bool
		}
		ValueAsString struct {
			Val string
			OK  bool
		}
		ValueAsBytes struct {
			Val []byte
			OK  bool
		}
		ValueAsInt struct {
			Val int
			OK  bool
		}
		ValueAsFloat struct {
			Val float64
			OK  bool
		}
		ValueAsBool struct {
			Val bool
			OK  bool
		}
		Variables struct {
			awake.Variables
		}
	}
}

func (s ScopeMock) ValueAsRaw(name string) (interface{}, bool) {
	s.Receives.ValueAsRaw.Name = name
	return s.Returns.ValueAsRaw.Val, s.Returns.ValueAsRaw.OK
}

func (s ScopeMock) ValueAsString(name string) (string, bool) {
	s.Receives.ValueAsString.Name = name
	return s.Returns.ValueAsString.Val, s.Returns.ValueAsString.OK
}

func (s ScopeMock) ValueAsBytes(name string) ([]byte, bool) {
	s.Receives.ValueAsBytes.Name = name
	return s.Returns.ValueAsBytes.Val, s.Returns.ValueAsBytes.OK
}

func (s ScopeMock) ValueAsInt(name string) (int, bool) {
	s.Receives.ValueAsInt.Name = name
	return s.Returns.ValueAsInt.Val, s.Returns.ValueAsInt.OK
}

func (s ScopeMock) ValueAsFloat(name string) (float64, bool) {
	s.Receives.ValueAsFloat.Name = name
	return s.Returns.ValueAsFloat.Val, s.Returns.ValueAsFloat.OK
}

func (s ScopeMock) ValueAsBool(name string) (bool, bool) {
	s.Receives.ValueAsBool.Name = name
	return s.Returns.ValueAsBool.Val, s.Returns.ValueAsBool.OK
}

func (s ScopeMock) SetReturnVariable(name string, variable awake.Variable) {
	s.Receives.SetReturnVariable.Name = name
	s.Receives.SetReturnVariable.Variable = variable
}

func (s ScopeMock) Variables() awake.Variables {
	return s.Returns.Variables.Variables
}
