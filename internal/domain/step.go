package domain

type Step struct {
	Conf   StepConfig `json:"conf"`
	Result StepResult `json:"result"`
	Err    error      `json:"err"`
	ErrMsg string     `json:"errMsg"`
}
