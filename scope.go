package awake

type Type string

const (
	TypeInt     Type = "integer"
	TypeString       = "string"
	TypeFloat        = "float"
	TypeBool         = "bool"
	TypeBytes        = "bytes"
	TypeDynamic      = "dynamic"
	TypeAny          = "any"
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
}
