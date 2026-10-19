package exercise

import (
	"fmt"
	"strings"
	"testing"
	"testing/fstest"
)

const modulesTOML = `
[[modules]]
id = "iteration"
title = "Iteration and Unpacking"
order = 1

[[modules]]
id = "hashing"
title = "Dictionaries and Sets"
order = 2

[[modules]]
id = "patterns"
title = "Applied Interview Patterns"
order = 7
`

func exerciseTOML(id, module string, order int) string {
	return fmt.Sprintf(`
id = %q
title = "Title %s"
module = %q
difficulty = "beginner"
type = "syntax"
order = %d
concepts = ["enumerate"]
description = "Do the thing."
objective = "Learn the thing."
interview = "Used everywhere."
hints = ["Try it."]

[[examples]]
input = "x"
output = "y"
`, id, id, module, order)
}

func addExercise(fsys fstest.MapFS, module, id string, order int) {
	dir := module + "/" + id + "/"
	fsys[dir+MetadataFile] = &fstest.MapFile{Data: []byte(exerciseTOML(id, module, order))}
	for _, name := range []string{StarterFile, TestFile, SolutionFile, ExplanationFile} {
		fsys[dir+name] = &fstest.MapFile{Data: []byte("# " + name)}
	}
}

func validFS() fstest.MapFS {
	fsys := fstest.MapFS{ModulesFile: {Data: []byte(modulesTOML)}}
	addExercise(fsys, "iteration", "zip-01", 2)
	addExercise(fsys, "iteration", "enumerate-01", 1)
	addExercise(fsys, "hashing", "sets-01", 1)
	return fsys
}

func TestLoadValid(t *testing.T) {
	cat, err := Load(validFS())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := len(cat.Modules); got != 2 {
		t.Fatalf("modules = %d, want 2 (empty module omitted)", got)
	}
	var ids []string
	for _, ex := range cat.All() {
		ids = append(ids, ex.ID)
	}
	if got, want := strings.Join(ids, ","), "enumerate-01,zip-01,sets-01"; got != want {
		t.Errorf("order = %s, want %s", got, want)
	}
	ex, ok := cat.Get("zip-01")
	if !ok {
		t.Fatal("Get(zip-01) not found")
	}
	if ex.TypeLabel() != "Syntax Drill" || len(ex.Examples) != 1 || ex.Dir != "iteration/zip-01" {
		t.Errorf("unexpected exercise: %+v", ex)
	}
	data, err := cat.ReadFile(ex, StarterFile)
	if err != nil || string(data) != "# starter.py" {
		t.Errorf("ReadFile = %q, %v", data, err)
	}
	if _, ok := cat.Get("missing"); ok {
		t.Error("Get(missing) found an exercise")
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name   string
		modify func(fstest.MapFS)
		want   string
	}{
		{
			name:   "missing modules file",
			modify: func(f fstest.MapFS) { delete(f, ModulesFile) },
			want:   "read modules.toml",
		},
		{
			name:   "invalid toml",
			modify: func(f fstest.MapFS) { f["iteration/zip-01/"+MetadataFile].Data = []byte("id = ") },
			want:   "parse exercise.toml",
		},
		{
			name: "duplicate id",
			modify: func(f fstest.MapFS) {
				addExercise(f, "hashing", "zip-01", 3)
			},
			want: `duplicate exercise id "zip-01"`,
		},
		{
			name:   "unlisted module directory",
			modify: func(f fstest.MapFS) { addExercise(f, "mystery", "x-01", 1) },
			want:   "mystery: directory is not listed",
		},
		{
			name:   "missing file",
			modify: func(f fstest.MapFS) { delete(f, "iteration/zip-01/"+SolutionFile) },
			want:   "missing file solution.py",
		},
		{
			name: "missing field and bad enums",
			modify: func(f fstest.MapFS) {
				f["iteration/zip-01/"+MetadataFile].Data = []byte(`
id = "zip-01"
module = "iteration"
difficulty = "easy"
type = "drill"
hints = ["h"]
`)
			},
			want: `missing "title"`,
		},
		{
			name: "unknown key",
			modify: func(f fstest.MapFS) {
				m := f["iteration/zip-01/"+MetadataFile]
				m.Data = append([]byte("hint = \"typo\"\n"), m.Data...)
			},
			want: "unknown keys: hint",
		},
		{
			name: "id does not match directory",
			modify: func(f fstest.MapFS) {
				f["iteration/zip-01/"+MetadataFile].Data = []byte(exerciseTOML("zip-02", "iteration", 1))
			},
			want: `must match directory name "zip-01"`,
		},
		{
			name: "module does not match directory",
			modify: func(f fstest.MapFS) {
				f["iteration/zip-01/"+MetadataFile].Data = []byte(exerciseTOML("zip-01", "hashing", 1))
			},
			want: `must match parent directory "iteration"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := validFS()
			tt.modify(fsys)
			_, err := Load(fsys)
			if err == nil {
				t.Fatalf("Load succeeded, want error containing %q", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %v\nwant substring %q", err, tt.want)
			}
		})
	}
}

func TestLoadReportsAllErrors(t *testing.T) {
	fsys := validFS()
	delete(fsys, "iteration/zip-01/"+SolutionFile)
	delete(fsys, "hashing/sets-01/"+TestFile)
	_, err := Load(fsys)
	if err == nil || !strings.Contains(err.Error(), "solution.py") || !strings.Contains(err.Error(), "test_exercise.py") {
		t.Errorf("want both problems reported, got %v", err)
	}
}

func TestLoadEmptyCurriculum(t *testing.T) {
	cat, err := Load(fstest.MapFS{ModulesFile: {Data: []byte(modulesTOML)}})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cat.All()) != 0 || len(cat.Modules) != 0 {
		t.Errorf("want empty catalog, got %d exercises", len(cat.All()))
	}
}
