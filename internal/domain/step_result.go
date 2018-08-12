package domain

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jvikstedt/awake"
)

type StepResult struct {
	Variables awake.Variables `json:"variables"`
}

type StepResults []StepResult

func (s *StepResults) Value() (driver.Value, error) {
	if s != nil {
		s, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		return string(s), nil
	}
	return nil, nil
}

func (s *StepResults) Scan(src interface{}) error {
	var data []byte
	if b, ok := src.([]byte); ok {
		data = b
	} else if s, ok := src.(string); ok {
		data = []byte(s)
	}
	return json.Unmarshal(data, s)
}
