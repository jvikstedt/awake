package awake

type Variable struct {
	Type string      `json:"type"`
	Val  interface{} `json:"val"`
}

type Step struct {
	Tag
}
