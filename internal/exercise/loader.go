package exercise

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// Catalog is the loaded curriculum. Modules without exercises are omitted.
type Catalog struct {
	FS      fs.FS
	Modules []*Module

	byID  map[string]*Exercise
	order []*Exercise
}

func (c *Catalog) Get(id string) (*Exercise, bool) {
	ex, ok := c.byID[id]
	return ex, ok
}

func (c *Catalog) All() []*Exercise { return c.order }

func (c *Catalog) ReadFile(ex *Exercise, name string) ([]byte, error) {
	return fs.ReadFile(c.FS, ex.File(name))
}

// Load reports every validation problem at once so authors can fix them together.
// Module and exercise directories are named NN_<id>; the number sets the order.
func Load(fsys fs.FS) (*Catalog, error) {
	modules, err := loadModules(fsys)
	if err != nil {
		return nil, err
	}
	moduleByID := make(map[string]*Module, len(modules))
	for _, m := range modules {
		moduleByID[m.ID] = m
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read exercise root: %w", err)
	}

	var errs []error
	cat := &Catalog{FS: fsys, byID: map[string]*Exercise{}}
	var found []*Module
	moduleNumbers := map[int]string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		num, id, ok := splitPrefix(entry.Name())
		if !ok {
			errs = append(errs, fmt.Errorf("%s: module directories must be named NN_<id>", entry.Name()))
			continue
		}
		mod, ok := moduleByID[id]
		if !ok {
			errs = append(errs, fmt.Errorf("%s: module %q is not listed in %s", entry.Name(), id, ModulesFile))
			continue
		}
		if prev, dup := moduleNumbers[num]; dup {
			errs = append(errs, fmt.Errorf("%s: number %02d is also used by %s", entry.Name(), num, prev))
			continue
		}
		moduleNumbers[num] = entry.Name()
		mod.Order = num
		found = append(found, mod)
		errs = append(errs, loadModuleExercises(fsys, cat, mod, entry.Name())...)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	slices.SortFunc(found, func(a, b *Module) int { return cmp.Compare(a.Order, b.Order) })
	for _, m := range found {
		if len(m.Exercises) == 0 {
			continue
		}
		slices.SortFunc(m.Exercises, func(a, b *Exercise) int { return cmp.Compare(a.Order, b.Order) })
		cat.Modules = append(cat.Modules, m)
		cat.order = append(cat.order, m.Exercises...)
	}
	return cat, nil
}

func loadModuleExercises(fsys fs.FS, cat *Catalog, mod *Module, dir string) []error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return []error{fmt.Errorf("%s: %w", dir, err)}
	}
	var errs []error
	numbers := map[int]string{}
	for _, d := range entries {
		if !d.IsDir() {
			continue
		}
		exDir := path.Join(dir, d.Name())
		num, id, ok := splitPrefix(d.Name())
		if !ok {
			errs = append(errs, fmt.Errorf("%s: exercise directories must be named NN_<id>", exDir))
			continue
		}
		if prev, dup := numbers[num]; dup {
			errs = append(errs, fmt.Errorf("%s: number %02d is also used by %s", exDir, num, prev))
			continue
		}
		numbers[num] = d.Name()
		ex, err := loadExercise(fsys, exDir, id, mod.ID)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if prev, dup := cat.byID[ex.ID]; dup {
			errs = append(errs, fmt.Errorf("%s: duplicate exercise id %q (also in %s)", ex.Dir, ex.ID, prev.Dir))
			continue
		}
		ex.Order = num
		cat.byID[ex.ID] = ex
		mod.Exercises = append(mod.Exercises, ex)
	}
	return errs
}

// splitPrefix parses directory names of the form "NN_name" with NN >= 01.
func splitPrefix(name string) (num int, rest string, ok bool) {
	digits, rest, found := strings.Cut(name, "_")
	if !found || len(digits) != 2 || rest == "" {
		return 0, "", false
	}
	num, err := strconv.Atoi(digits)
	if err != nil || num < 1 {
		return 0, "", false
	}
	return num, rest, true
}

func loadModules(fsys fs.FS) ([]*Module, error) {
	data, err := fs.ReadFile(fsys, ModulesFile)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", ModulesFile, err)
	}
	var doc struct {
		Modules []*Module `toml:"modules"`
	}
	md, err := toml.Decode(string(data), &doc)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", ModulesFile, err)
	}
	if err := checkUndecoded(md); err != nil {
		return nil, fmt.Errorf("%s: %w", ModulesFile, err)
	}

	seen := map[string]bool{}
	var errs []error
	for i, m := range doc.Modules {
		switch {
		case m.ID == "" || m.Title == "":
			errs = append(errs, fmt.Errorf("%s: module %d needs both id and title", ModulesFile, i+1))
		case seen[m.ID]:
			errs = append(errs, fmt.Errorf("%s: duplicate module id %q", ModulesFile, m.ID))
		}
		seen[m.ID] = true
	}
	if len(doc.Modules) == 0 {
		errs = append(errs, fmt.Errorf("%s: no modules defined", ModulesFile))
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return doc.Modules, nil
}

func loadExercise(fsys fs.FS, dir, id, module string) (*Exercise, error) {
	data, err := fs.ReadFile(fsys, path.Join(dir, MetadataFile))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dir, err)
	}
	ex := &Exercise{Dir: dir}
	md, err := toml.Decode(string(data), ex)
	if err != nil {
		return nil, fmt.Errorf("%s: parse %s: %w", dir, MetadataFile, err)
	}
	if err := checkUndecoded(md); err != nil {
		return nil, fmt.Errorf("%s: %w", dir, err)
	}
	if err := validate(fsys, ex, id, module); err != nil {
		return nil, fmt.Errorf("%s: %w", dir, err)
	}
	return ex, nil
}

// checkUndecoded rejects unknown keys so typos like "hint =" fail loudly.
func checkUndecoded(md toml.MetaData) error {
	keys := md.Undecoded()
	if len(keys) == 0 {
		return nil
	}
	names := make([]string, len(keys))
	for i, k := range keys {
		names[i] = k.String()
	}
	return fmt.Errorf("unknown keys: %s", strings.Join(names, ", "))
}
