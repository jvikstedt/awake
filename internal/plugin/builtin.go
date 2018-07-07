package plugin

import (
	"github.com/jvikstedt/awake/internal/domain"
	"github.com/jvikstedt/awake/internal/plugin/builtin"
)

func BuiltinPerformers() []domain.Performer {
	performers := []domain.Performer{
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
