package exercise

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

func validate(fsys fs.FS, ex *Exercise, dirName, module string) error {
	var errs []error
	required := []struct{ name, value string }{
		{"id", ex.ID},
		{"title", ex.Title},
		{"description", ex.Description},
		{"objective", ex.Objective},
		{"interview", ex.Interview},
	}
	for _, f := range required {
		if strings.TrimSpace(f.value) == "" {
			errs = append(errs, fmt.Errorf("missing %q", f.name))
		}
	}
	if ex.ID != "" && ex.ID != dirName {
		errs = append(errs, fmt.Errorf("id %q must match directory name %q", ex.ID, dirName))
	}
	if ex.Module != module {
		errs = append(errs, fmt.Errorf("module %q must match parent directory %q", ex.Module, module))
	}
	if !validDifficulties[ex.Difficulty] {
		errs = append(errs, fmt.Errorf("difficulty %q must be beginner, intermediate, or advanced", ex.Difficulty))
	}
	if _, ok := validTypes[ex.Type]; !ok {
		errs = append(errs, fmt.Errorf("type %q must be %q or %q", ex.Type, TypeSyntax, TypeApplication))
	}
	if len(ex.Hints) == 0 {
		errs = append(errs, errors.New("needs at least one hint"))
	}
	if len(ex.Examples) == 0 {
		errs = append(errs, errors.New("needs at least one example"))
	}
	for i, e := range ex.Examples {
		if strings.TrimSpace(e.Input) == "" || strings.TrimSpace(e.Output) == "" || strings.TrimSpace(e.Explanation) == "" {
			errs = append(errs, fmt.Errorf("example %d needs input, output, and explanation", i+1))
		}
	}
	for _, name := range requiredFiles {
		if _, err := fs.Stat(fsys, ex.File(name)); err != nil {
			errs = append(errs, fmt.Errorf("missing file %s", name))
		}
	}
	return errors.Join(errs...)
}
