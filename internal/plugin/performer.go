package plugin

import "github.com/jvikstedt/awake"

// Performer is registered to system as plugin and are used to execute steps.
type Performer interface {
	Info() awake.PerformerInfo
	Perform(awake.Scope) error
}

// Tag is just a string that defines name of performer
type Tag string

var performers = map[Tag]Performer{}

// RegisterPerformer registers performer by tag to global variable
func RegisterPerformer(p Performer) {
	performers[Tag(p.Info().Name)] = p
}

func FindPerformer(tag Tag) (Performer, bool) {
	i, ok := performers[tag]
	return i, ok
}
