package builtin

import (
	"github.com/jvikstedt/awake"
)

type Alerts struct{}

func (Alerts) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_alerts",
		DisplayName: "Alert count",
	}
}

func (Alerts) Perform(scope awake.Scope) error {
	alerts := scope.Alerts()

	errors := alerts.ByType(awake.AlertError)
	warnings := alerts.ByType(awake.AlertWarning)

	scope.SetReturnVariable("errorCount", awake.Variable{Type: awake.TypeInt, Val: len(errors)})
	scope.SetReturnVariable("hasErrors", awake.Variable{Type: awake.TypeBool, Val: len(errors) > 0})

	scope.SetReturnVariable("warningCount", awake.Variable{Type: awake.TypeInt, Val: len(warnings)})
	scope.SetReturnVariable("hasWarnings", awake.Variable{Type: awake.TypeBool, Val: len(warnings) > 0})

	alertStr := ""
	for _, a := range alerts {
		alertStr += a.Value + "\n"
	}

	scope.SetReturnVariable("alerts", awake.Variable{Type: awake.TypeString, Val: alertStr})

	return nil
}
