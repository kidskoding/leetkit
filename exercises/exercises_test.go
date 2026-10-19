package exercises_test

import (
	"testing"

	"github.com/anirudhkonidala/leetkit/exercises"
	"github.com/anirudhkonidala/leetkit/internal/exercise"
)

func TestBundledCurriculumLoads(t *testing.T) {
	cat, err := exercise.Load(exercises.FS)
	if err != nil {
		t.Fatalf("bundled curriculum is invalid:\n%v", err)
	}
	if len(cat.All()) == 0 {
		t.Fatal("bundled curriculum has no exercises")
	}
}
