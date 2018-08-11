package awake

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

type PerformerInfo struct {
	Name        string
	DisplayName string
}

type Type string

const (
	TypeInt        Type = "integer"
	TypeString          = "string"
	TypeFloat           = "float"
	TypeBool            = "bool"
	TypeByte            = "byte"
	TypeDynamic         = "dynamic"
	TypeArrayBytes      = "bytes"
	TypeNil             = "nil"
	TypeAny             = "any"
)

type Variable struct {
	Type `json:"type"`
	Val  interface{} `json:"val"`
}

type Variables map[string]Variable

// Scope is passed to performers and defines how performer can access data
type Scope interface {
	ValueAsRaw(name string) (interface{}, bool)
	ValueAsString(name string) (string, bool)
	ValueAsBytes(name string) ([]byte, bool)
	ValueAsInt(name string) (int, bool)
	ValueAsFloat(name string) (float64, bool)
	ValueAsBool(name string) (bool, bool)
	SetReturnVariable(name string, variable Variable)
	Variables() Variables
	Errors() []error
}

func ResolveType(i interface{}) Type {
	if i == nil {
		return TypeNil
	}
	switch i.(type) {
	case int:
		return TypeInt
	case string:
		return TypeString
	case float64:
		return TypeFloat
	case float32:
		return TypeFloat
	case bool:
		return TypeBool
	}

	if reflect.TypeOf(i).Kind() == reflect.Slice {
		if reflect.TypeOf(i).Elem().Kind() == reflect.Uint8 {
			return TypeArrayBytes
		}
	}

	return TypeAny
}

func MakeVariable(i interface{}) Variable {
	return Variable{Type: ResolveType(i), Val: i}
}

func (v Variables) Value() (driver.Value, error) {
	if v != nil {
		v, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return string(v), nil
	}
	return nil, nil
}

func (v Variables) Scan(src interface{}) error {
	var data []byte
	if b, ok := src.([]byte); ok {
		data = b
	} else if v, ok := src.(string); ok {
		data = []byte(v)
	}
	return json.Unmarshal(data, &v)
}
