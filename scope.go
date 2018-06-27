package awake

type Variable struct {
	Type string      `json:"type"`
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
}
