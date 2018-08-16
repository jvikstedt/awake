package runner

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/domain"
)

var logs = &bytes.Buffer{}
var logger = log.New(logs, "", log.LstdFlags)

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}

func TestValueAsRaw(t *testing.T) {
	tt := []struct {
		tname string
		name  string
		val   interface{}
		ok    bool
	}{
		{tname: "int", name: "code", val: 200, ok: true},
		{tname: "float64", name: "temp", val: 21.2, ok: true},
		{tname: "float64 to int", name: "code2", val: 200, ok: true},
		{tname: "string", name: "status", val: "foo", ok: true},
		{tname: "missing", name: "foo", ok: false},
	}

	variables := awake.Variables{
		"code": awake.Variable{
			Type: awake.TypeInt,
			Val:  200,
		},
		"code2": awake.Variable{
			Type: awake.TypeInt,
			Val:  200.0,
		},
		"status": awake.Variable{
			Type: awake.TypeString,
			Val:  "foo",
		},
		"temp": awake.Variable{
			Type: awake.TypeFloat,
			Val:  21.2,
		},
	}

	stepConfigs := []domain.StepConfig{
		{Tag: "FOO", Variables: variables},
	}

	scope := newScope(logger, domain.PerformerConfigs{}, stepConfigs, domain.Storage{})
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			val, ok := scope.ValueAsRaw(tc.name)
			if ok != tc.ok {
				t.Fatalf("Expected %v to eq %v", ok, tc.ok)
			}
			if val != tc.val {
				t.Fatalf("Expected %v to eq %v", val, tc.val)
			}
		})
	}
}

func TestValueAsRawDynamic(t *testing.T) {
	tt := []struct {
		tname string
		name  string
		val   interface{}
		ok    bool
	}{
		{tname: "int variable", name: "code", val: 200.0, ok: true},
		{tname: "one variable in string", name: "header", val: "authorization: Bearer abc123", ok: true},
		{tname: "two variable in string", name: "credentials", val: `{"username":"foo","password":"bar"}`, ok: true},
		{tname: "missing", name: "foo", val: nil, ok: false},
	}

	steps := []domain.StepConfig{
		{Tag: "HTTP", Variables: awake.Variables{}},
		{
			Tag: "ASSERT",
			Variables: awake.Variables{
				"header": awake.Variable{
					Type: awake.TypeDynamic,
					Val:  "'authorization: Bearer ' + [0-token]",
				},
				"code": awake.Variable{
					Type: awake.TypeDynamic,
					Val:  "[0-code]",
				},
				"credentials": awake.Variable{
					Type: awake.TypeDynamic,
					Val:  `'{\"username\":\"' + [0-username] + '\",\"password\":\"' + [0-password] + '\"}'`,
				},
			},
		},
	}

	scope := newScope(logger, domain.PerformerConfigs{}, steps, domain.Storage{})
	scope.current = 1

	scope.steps[0].Result.Variables["token"] = awake.Variable{Type: awake.TypeString, Val: "abc123"}
	scope.steps[0].Result.Variables["code"] = awake.Variable{Type: awake.TypeInt, Val: 200.0}
	scope.steps[0].Result.Variables["username"] = awake.Variable{Type: awake.TypeString, Val: "foo"}
	scope.steps[0].Result.Variables["password"] = awake.Variable{Type: awake.TypeString, Val: "bar"}

	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			val, ok := scope.ValueAsRaw(tc.name)
			if ok != tc.ok {
				t.Fatalf("Expected %v to eq %v", ok, tc.ok)
			}
			if val != tc.val {
				t.Fatalf("Expected %v to eq %v", val, tc.val)
			}

		})
	}
}
