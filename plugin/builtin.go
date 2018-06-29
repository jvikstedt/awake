package plugin

import (
	"github.com/jvikstedt/awake/internal/task"
	"github.com/jvikstedt/awake/plugin/builtin"
)

func BuiltinPerformers() []task.Performer {
	performers := []task.Performer{
		builtin.Equal{},
		builtin.HTTP{},
		builtin.JSON{},
		builtin.DiskUsage{},
		builtin.LessThan{},
		builtin.MoreThan{},
	}

	return performers
}
