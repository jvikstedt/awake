package domain

import (
	"database/sql/driver"
	"encoding/json"
)

type Step struct {
	Conf   StepConfig `json:"conf"`
	Result StepResult `json:"result"`
	Err    error      `json:"err"`
	ErrMsg string     `json:"errMsg"`
}

type Steps []Step

func (s *Steps) Value() (driver.Value, error) {
	if s != nil {
		s, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		return string(s), nil
	}
	return nil, nil
}

func (s *Steps) Scan(src interface{}) error {
	var data []byte
	if b, ok := src.([]byte); ok {
		data = b
	} else if s, ok := src.(string); ok {
		data = []byte(s)
	}
	return json.Unmarshal(data, s)
}
