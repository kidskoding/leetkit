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

[[modules]]
id = "hashing"
title = "Dictionaries and Sets"

[[modules]]
id = "patterns"
title = "Applied Interview Patterns"
`

func exerciseTOML(id, module string) string {
	return fmt.Sprintf(`
id = %q
title = "Title %s"
module = %q
difficulty = "beginner"
type = "syntax"
concepts = ["enumerate"]
description = "Do the thing."
objective = "Learn the thing."
interview = "Used everywhere."
hints = ["Try it."]

[[examples]]
input = "x"
output = "y"
explanation = "Because."
`, id, id, module)
}

var moduleDirs = map[string]string{"iteration": "01_iteration", "hashing": "02_hashing", "patterns": "07_patterns", "mystery": "09_mystery"}

// exDir is the directory of exercise id, numbered n within its module.
func exDir(module, id string, n int) string {
	return fmt.Sprintf("%s/%02d_%s/", moduleDirs[module], n, id)
}

func addExercise(fsys fstest.MapFS, module, id string, order int) {
	dir := exDir(module, id, order)
	fsys[dir+MetadataFile] = &fstest.MapFile{Data: []byte(exerciseTOML(id, module))}
	for _, name := range []string{StarterFile, TestFile, SolutionFile, ExplanationFile} {
		fsys[dir+name] = &fstest.MapFile{Data: []byte("# " + name)}
	}
}

var zipDir = exDir("iteration", "zip-01", 2)

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
	if ex.TypeLabel() != "Syntax Drill" || len(ex.Examples) != 1 || ex.Dir != "01_iteration/02_zip-01" {
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
			modify: func(f fstest.MapFS) { f[zipDir+MetadataFile].Data = []byte("id = ") },
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
			want:   `09_mystery: module "mystery" is not listed`,
		},
		{
			name:   "unnumbered exercise directory",
			modify: func(f fstest.MapFS) { f["01_iteration/zip-02/"+MetadataFile] = &fstest.MapFile{Data: []byte("")} },
			want:   "exercise directories must be named NN_<id>",
		},
		{
			name:   "unnumbered module directory",
			modify: func(f fstest.MapFS) { f["iteration/README.md"] = &fstest.MapFile{} },
			want:   "module directories must be named NN_<id>",
		},
		{
			name:   "duplicate exercise number",
			modify: func(f fstest.MapFS) { addExercise(f, "iteration", "range-01", 2) },
			want:   "number 02 is also used by",
		},
		{
			name: "order key is rejected",
			modify: func(f fstest.MapFS) {
				m := f[zipDir+MetadataFile]
				m.Data = append([]byte("order = 2\n"), m.Data...)
			},
			want: "unknown keys: order",
		},
		{
			name:   "missing file",
			modify: func(f fstest.MapFS) { delete(f, zipDir+SolutionFile) },
			want:   "missing file solution.py",
		},
		{
			name: "missing field and bad enums",
			modify: func(f fstest.MapFS) {
				f[zipDir+MetadataFile].Data = []byte(`
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
				m := f[zipDir+MetadataFile]
				m.Data = append([]byte("hint = \"typo\"\n"), m.Data...)
			},
			want: "unknown keys: hint",
		},
		{
			name: "id does not match directory",
			modify: func(f fstest.MapFS) {
				f[zipDir+MetadataFile].Data = []byte(exerciseTOML("zip-02", "iteration"))
			},
			want: `must match directory name "zip-01"`,
		},
		{
			name: "module does not match directory",
			modify: func(f fstest.MapFS) {
				f[zipDir+MetadataFile].Data = []byte(exerciseTOML("zip-01", "hashing"))
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
	delete(fsys, zipDir+SolutionFile)
	delete(fsys, exDir("hashing", "sets-01", 1)+TestFile)
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
