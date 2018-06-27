package plugin

import (
	"github.com/jvikstedt/awake/internal/task"
	"github.com/jvikstedt/awake/plugin/builtin"
)

func BuiltinPerformers(fn func(task.Performer)) {
	performers := []task.Performer{
		builtin.Equal{},
		builtin.HTTP{},
		builtin.JSON{},
	}

	for _, v := range performers {
		fn(v)
	}
}
