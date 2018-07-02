package domain

type Step struct {
	Conf   StepConfig
	Result StepResult
	Err    error
	ErrMsg string
}
