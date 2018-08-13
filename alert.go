package awake

type AlertType string

const (
	AlertError   AlertType = "error"
	AlertWarning AlertType = "warning"
)

type Alert struct {
	Type  AlertType `json:"type"`
	Value string    `json:"value"`
	Meta  Variables `json:"meta"`
}

type Alerts []Alert

func (a Alerts) ByType(t AlertType) Alerts {
	alerts := Alerts{}

	for _, v := range a {
		if v.Type == t {
			alerts = append(alerts, v)
		}
	}

	return alerts
}
