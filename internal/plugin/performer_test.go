package plugin_test

import (
	"testing"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/plugin"
)

type testPerformer struct {
	Name string
}

func (tp testPerformer) Info() awake.PerformerInfo {
	return awake.PerformerInfo{Name: tp.Name, DisplayName: "test"}
}

func (tp testPerformer) Perform(awake.Scope) error { return nil }

func TestRegisterPerformer(t *testing.T) {
	tt := []struct {
		tag string
	}{
		{tag: "HTTP"},
		{tag: "EQUAL"},
	}

	for _, v := range tt {
		t.Run(v.tag, func(t *testing.T) {
			tp := testPerformer{Name: v.tag}
			plugin.RegisterPerformer(tp)
			rp, ok := plugin.FindPerformer(plugin.Tag(v.tag))
			if !ok {
				t.Fatalf("Could not find performer by tag %s", v.tag)
			}

			if tp != rp {
				t.Fatalf("Expected %v but got %v", tp, rp)
			}
		})
	}

	t.Run("NOT_FOUND", func(t *testing.T) {
		rp, ok := plugin.FindPerformer("FOO")
		if ok {
			t.Fatalf("Expected not ok but was ok")
		}

		if rp != nil {
			t.Fatalf("Expected performer to be nil but was %v", rp)
		}
	})
}
