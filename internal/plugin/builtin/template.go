package builtin

import (
	"bytes"
	templ "text/template"

	"github.com/jvikstedt/awake"
)

type Template struct{}

func (Template) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_template",
		DisplayName: "Template",
	}
}

func (Template) Perform(scope awake.Scope) error {
	template, _ := scope.ValueAsString("template")
	variables := scope.Variables()

	tmpl, err := templ.New("template").Parse(template)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{}
	for key, v := range variables {
		vars[key] = v.Val
	}

	buffer := bytes.Buffer{}
	err = tmpl.Execute(&buffer, vars)
	if err != nil {
		return err
	}

	scope.SetReturnVariable("return", awake.Variable{Type: awake.TypeString, Val: buffer.String()})

	return nil
}
