package domain

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/plugin"
)

type StepConfig struct {
	plugin.Tag      `json:"tag"`
	awake.Variables `json:"variables"`
}

type StepConfigs []StepConfig

func (s *StepConfigs) Value() (driver.Value, error) {
	if s != nil {
		s, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		return string(s), nil
	}
	return nil, nil
}

func (s *StepConfigs) Scan(src interface{}) error {
	var data []byte
	if b, ok := src.([]byte); ok {
		data = b
	} else if s, ok := src.(string); ok {
		data = []byte(s)
	}
	return json.Unmarshal(data, s)
}
