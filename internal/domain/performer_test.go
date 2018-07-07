package domain_test

import (
	"testing"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/domain"
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
		tag domain.Tag
	}{
		{tag: "HTTP"},
		{tag: "EQUAL"},
	}

	for _, v := range tt {
		t.Run(string(v.tag), func(t *testing.T) {
			tp := testPerformer{Name: string(v.tag)}
			domain.RegisterPerformer(tp)
			rp, ok := domain.FindPerformer(v.tag)
			if !ok {
				t.Fatalf("Could not find performer by tag %s", v.tag)
			}

			if tp != rp {
				t.Fatalf("Expected %v but got %v", tp, rp)
			}
		})
	}

	t.Run("NOT_FOUND", func(t *testing.T) {
		rp, ok := domain.FindPerformer("FOO")
		if ok {
			t.Fatalf("Expected not ok but was ok")
		}

		if rp != nil {
			t.Fatalf("Expected performer to be nil but was %v", rp)
		}
	})
}
