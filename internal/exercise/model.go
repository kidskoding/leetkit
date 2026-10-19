// Package exercise loads and validates the curriculum.
package exercise

import "path"

const (
	MetadataFile    = "exercise.toml"
	StarterFile     = "starter.py"
	TestFile        = "test_exercise.py"
	SolutionFile    = "solution.py"
	ExplanationFile = "explanation.md"
)

const ModulesFile = "modules.toml"

var requiredFiles = []string{MetadataFile, StarterFile, TestFile, SolutionFile, ExplanationFile}

const (
	TypeSyntax      = "syntax"
	TypeApplication = "application"
)

var validTypes = map[string]string{
	TypeSyntax:      "Syntax Drill",
	TypeApplication: "Interview Application",
}

var validDifficulties = map[string]bool{
	"beginner":     true,
	"intermediate": true,
	"advanced":     true,
}

type Example struct {
	Input       string `toml:"input"`
	Output      string `toml:"output"`
	Explanation string `toml:"explanation"`
}

type Exercise struct {
	ID         string `toml:"id"`
	Title      string `toml:"title"`
	Module     string `toml:"module"`
	Difficulty string `toml:"difficulty"`
	Type       string `toml:"type"`
	// Order comes from the directory's NN_ prefix.
	Order       int       `toml:"-"`
	Concepts    []string  `toml:"concepts"`
	Description string    `toml:"description"`
	Objective   string    `toml:"objective"`
	Interview   string    `toml:"interview"`
	Hints       []string  `toml:"hints"`
	Examples    []Example `toml:"examples"`

	Dir string `toml:"-"`
}

func (e *Exercise) TypeLabel() string { return validTypes[e.Type] }

func (e *Exercise) File(name string) string { return path.Join(e.Dir, name) }

type Module struct {
	ID    string `toml:"id"`
	Title string `toml:"title"`
	// Order comes from the directory's NN_ prefix.
	Order     int         `toml:"-"`
	Exercises []*Exercise `toml:"-"`
}
