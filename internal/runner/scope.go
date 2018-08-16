package runner

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/Knetic/govaluate"
	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/domain"
)

var dynamicRegexp = regexp.MustCompile(`\${[^}]*}`)

type scope struct {
	log     *log.Logger
	current int
	steps   domain.Steps
	storage domain.Storage
}

func newScope(l *log.Logger, performerConfigs domain.PerformerConfigs, stepConfigs []domain.StepConfig, storage domain.Storage) *scope {
	steps := make(domain.Steps, len(stepConfigs))

	for i, stepConfig := range stepConfigs {
		steps[i] = domain.Step{
			Conf: stepConfig,
			Result: domain.StepResult{
				Variables: awake.Variables{},
			},
			Alerts: awake.Alerts{},
		}

		conf, ok := performerConfigs[steps[i].Conf.Tag]
		if ok {
			for key, variable := range conf {
				if _, ok := steps[i].Conf.Variables[key]; !ok {
					steps[i].Conf.Variables[key] = variable
				}
			}
		}
	}

	return &scope{
		log:     l,
		current: 0,
		steps:   steps,
		storage: storage,
	}
}

func (s *scope) currentStep() domain.Step {
	return s.steps[s.current]
}

func (s *scope) getValue(v awake.Variable) (interface{}, error) {
	switch v.Type {
	case awake.TypeInt:
		return s.handleInt(v.Val)
	case awake.TypeDynamic:
		return s.handleDynamic(v.Val)
	default:
		return v.Val, nil
	}
}

func (s *scope) addAlert(alert awake.Alert) {
	s.steps[s.current].Alerts = append(s.currentStep().Alerts, alert)
}

// Implements awake.Scope

func (s *scope) SetReturnVariable(name string, variable awake.Variable) {
	s.currentStep().Result.Variables[name] = variable
}

func (s *scope) SetStorageVariable(name string, variable awake.Variable) {
	s.storage[name] = variable
}

func (s *scope) Alerts() awake.Alerts {
	alerts := awake.Alerts{}
	for i := 0; i < s.current; i++ {
		alerts = append(alerts, s.steps[i].Alerts...)
	}

	return alerts
}

func (s *scope) Variables() awake.Variables {
	vars := awake.Variables{}
	for key, v := range s.currentStep().Conf.Variables {
		val, err := s.getValue(v)
		if err != nil {
			s.log.Println(err)
			continue
		}
		if v.Type == awake.TypeDynamic {
			vars[key] = awake.MakeVariable(val)
			continue
		}
		vars[key] = awake.Variable{Type: v.Type, Val: val}
	}

	return vars
}

func (s *scope) ValueAsRaw(name string) (interface{}, bool) {
	currentStepPair := s.currentStep()

	v, ok := currentStepPair.Conf.Variables[name]
	if !ok {
		s.addAlert(awake.Alert{Type: awake.AlertWarning, Value: fmt.Sprintf("Could not find variable by name: %s\n", name)})
		return nil, ok
	}

	val, err := s.getValue(v)
	if err != nil {
		s.addAlert(awake.Alert{Type: awake.AlertWarning, Value: err.Error()})
		return nil, false
	}

	return val, true
}

func (s *scope) ValueAsString(name string) (string, bool) {
	v, ok := s.ValueAsRaw(name)
	if !ok {
		return "", ok
	}

	asStr, ok := v.(string)
	return asStr, ok
}

func (s *scope) ValueAsBytes(name string) ([]byte, bool) {
	v, ok := s.ValueAsRaw(name)
	if !ok {
		return []byte{}, ok
	}

	asBytes, ok := v.([]byte)
	return asBytes, ok
}

func (s *scope) ValueAsInt(name string) (int, bool) {
	v, ok := s.ValueAsRaw(name)
	if !ok {
		return 0, ok
	}

	val, err := s.handleInt(v)
	if err != nil {
		return 0, false
	}

	return val, true
}

func (s *scope) ValueAsFloat(name string) (float64, bool) {
	v, ok := s.ValueAsRaw(name)
	if !ok {
		return 0, ok
	}

	asFloat64, ok := v.(float64)
	return asFloat64, ok
}

func (s *scope) ValueAsBool(name string) (bool, bool) {
	v, ok := s.ValueAsRaw(name)
	if !ok {
		return false, ok
	}

	asBool, ok := v.(bool)
	return asBool, ok
}

func (s *scope) handleInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("Expected %v to be either int or float64, but was %T", val, val)
	}
}

var functions = map[string]govaluate.ExpressionFunction{
	"toInt": func(args ...interface{}) (interface{}, error) {
		val := args[0]

		switch v := val.(type) {
		case int:
			return v, nil
		case float64:
			return int(v), nil
		case string:
			return strconv.Atoi(v)
		default:
			return val, fmt.Errorf("Expected %v to be either int, string or float64, but was %T", val, val)
		}
	},
	"toFloat": func(args ...interface{}) (interface{}, error) {
		val := args[0]

		switch v := val.(type) {
		case int:
			return float64(v), nil
		case float64:
			return v, nil
		case string:
			return strconv.ParseFloat(v, 64)
		default:
			return val, fmt.Errorf("Expected %v to be either int or float64, but was %T", val, val)
		}
	},
}

func (s *scope) handleDynamic(val interface{}) (interface{}, error) {
	vars := map[string]interface{}{}
	for i, step := range s.steps {
		for key, v := range step.Result.Variables {
			vars[fmt.Sprintf("%d-%s", i, key)] = v.Val
		}
	}

	for key, val := range s.storage {
		vars[fmt.Sprintf("storage-%s", key)] = val.Val
	}

	asStr, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("Expected %v to be string but got %T", val, val)
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(asStr, functions)
	if err != nil {
		return nil, err
	}

	return expression.Evaluate(vars)
}
