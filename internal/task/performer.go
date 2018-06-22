package task

import "github.com/jvikstedt/awake"

// Performer is registered to system as plugin and are used to execute steps.
type Performer interface {
	Perform(awake.Scope) error
}

// Tag is just a string that defines name of performer
type Tag string

var performers = map[Tag]Performer{}

// RegisterPerformer registers performer by tag to global variable
func RegisterPerformer(tag Tag, p Performer) {
	performers[tag] = p
}

func FindPerformer(tag Tag) (Performer, bool) {
	i, ok := performers[tag]
	return i, ok
}
