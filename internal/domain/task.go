package domain

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/jvikstedt/awake"
)

var dynamicRegexp = regexp.MustCompile(`\${[^}]*}`)

type Task struct {
	log              *log.Logger
	current          int
	steps            []*Step
	performerConfigs PerformerConfigs
}

func NewTask(l *log.Logger, performerConfigs PerformerConfigs, stepConfigs []StepConfig) *Task {
	steps := make([]*Step, len(stepConfigs))

	for i, stepConfig := range stepConfigs {
		steps[i] = &Step{
			Conf: stepConfig,
			Result: StepResult{
				Variables: awake.Variables{},
			},
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

	return &Task{
		log:              l,
		current:          0,
		steps:            steps,
		performerConfigs: performerConfigs,
	}
}

func (t *Task) Run() []*Step {
	for i, v := range t.steps {
		t.current = i

		performer, ok := FindPerformer(v.Conf.Tag)
		if !ok {
			v.Err = fmt.Errorf("argh... performer not found %s", v.Conf.Tag)
			v.ErrMsg = v.Err.Error()
			continue
		}

		if v.Err = performer.Perform(t); v.Err != nil {
			v.ErrMsg = v.Err.Error()
		}
	}

	return t.steps
}

func (t *Task) SetReturnVariable(name string, variable awake.Variable) {
	t.currentStep().Result.Variables[name] = variable
}

func (t *Task) Errors() []error {
	errors := []error{}
	for i := 0; i < t.current; i++ {
		if t.steps[i].Err != nil {
			errors = append(errors, t.steps[i].Err)
		}
	}

	return errors
}

func (t *Task) Variables() awake.Variables {
	vars := awake.Variables{}
	for key, v := range t.currentStep().Conf.Variables {
		val, err := t.getValue(v)
		if err != nil {
			t.log.Println(err)
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

func (t *Task) currentStep() *Step {
	return t.steps[t.current]
}

// Implements awake.Scope

func (t *Task) ValueAsRaw(name string) (interface{}, bool) {
	currentStepPair := t.currentStep()

	v, ok := currentStepPair.Conf.Variables[name]
	if !ok {
		t.log.Printf("Could not find variable by name: %s\n", name)
		return nil, ok
	}

	val, err := t.getValue(v)
	if err != nil {
		t.log.Println(err)
		return nil, false
	}

	return val, true
}

func (t *Task) ValueAsString(name string) (string, bool) {
	v, ok := t.ValueAsRaw(name)
	if !ok {
		return "", ok
	}

	asStr, ok := v.(string)
	return asStr, ok
}

func (t *Task) ValueAsBytes(name string) ([]byte, bool) {
	v, ok := t.ValueAsRaw(name)
	if !ok {
		return []byte{}, ok
	}

	asBytes, ok := v.([]byte)
	return asBytes, ok
}

func (t *Task) ValueAsInt(name string) (int, bool) {
	v, ok := t.ValueAsRaw(name)
	if !ok {
		return 0, ok
	}

	val, err := t.handleInt(v)
	if err != nil {
		return 0, false
	}

	return val, true
}

func (t *Task) ValueAsFloat(name string) (float64, bool) {
	v, ok := t.ValueAsRaw(name)
	if !ok {
		return 0, ok
	}

	asFloat64, ok := v.(float64)
	return asFloat64, ok
}

func (t *Task) ValueAsBool(name string) (bool, bool) {
	v, ok := t.ValueAsRaw(name)
	if !ok {
		return false, ok
	}

	asBool, ok := v.(bool)
	return asBool, ok
}

func (t *Task) handleInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("Expected %v to be either int or float64, but was %T", val, val)
	}
}

func (t *Task) handleDynamic(val interface{}) (interface{}, error) {
	asStr, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("Expected %v to be string but got %T", val, val)
	}

	matches := dynamicRegexp.FindAllString(asStr, -1)

	if len(matches) == 1 && len(matches[0]) == len(asStr) {
		s := strings.Split(matches[0][2:len(matches[0])-1], ":")
		i, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, err
		}
		return t.steps[i].Result.Variables[s[1]].Val, nil
	}

	compiled := asStr

	for _, m := range matches {
		s := strings.Split(m[2:len(m)-1], ":")
		i, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, err
		}
		scopeValAsStr := fmt.Sprintf("%v", t.steps[i].Result.Variables[s[1]].Val)
		compiled = strings.Replace(compiled, m, scopeValAsStr, -1)
	}

	return compiled, nil
}

func (t *Task) getValue(v awake.Variable) (interface{}, error) {
	switch v.Type {
	case awake.TypeInt:
		return t.handleInt(v.Val)
	case awake.TypeDynamic:
		return t.handleDynamic(v.Val)
	default:
		return v.Val, nil
	}
}
