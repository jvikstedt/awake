package runner

import "github.com/jvikstedt/awake/internal/domain"

type Step struct {
	Conf   domain.StepConfig `json:"conf"`
	Result domain.StepResult `json:"result"`
	Err    error             `json:"err"`
	ErrMsg string            `json:"errMsg"`
}
