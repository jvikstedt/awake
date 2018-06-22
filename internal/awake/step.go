package awake

type Variable struct {
	Type string      `json:"type"`
	Val  interface{} `json:"val"`
}

type Variables map[string]Variable

type Step struct {
	Tag
	Variables
}
