package awake

// Scope is passed to performers and defines how performer can access data
type Scope interface {
	ValueAsRaw(name string) (interface{}, bool)
	ValueAsString(name string) (string, bool)
	ValueAsInt(name string) (int, bool)
	ValueAsFloat(name string) (float64, bool)
	ValueAsBool(name string) (bool, bool)
	SetReturnValue(name string, typ string, val interface{})
}
