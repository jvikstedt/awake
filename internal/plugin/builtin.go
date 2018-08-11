package plugin

import (
	"github.com/jvikstedt/awake/internal/plugin/builtin"
)

func BuiltinPerformers() []Performer {
	performers := []Performer{
		builtin.Equal{},
		builtin.HTTP{},
		builtin.JSON{},
		builtin.DiskUsage{},
		builtin.LessThan{},
		builtin.MoreThan{},
		builtin.Template{},
		builtin.Mailer{},
		builtin.Errors{},
	}

	return performers
}
