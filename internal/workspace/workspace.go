// Package workspace manages learners' working copies of exercise code. The
// bundled starter files are templates and are never modified.
package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kidskoding/leetkit/internal/atomicfile"
	"github.com/kidskoding/leetkit/internal/exercise"
)

type Workspace struct {
	Dir string
}

var errNoDir = errors.New("workspace directory is not set")

func (w Workspace) Path(id string) string {
	return filepath.Join(w.Dir, id, exercise.StarterFile)
}

// Ensure returns the working copy for id, creating it from starter on first use.
func (w Workspace) Ensure(id string, starter []byte) (string, error) {
	if w.Dir == "" {
		return "", errNoDir
	}
	path := w.Path(id)
	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("check working copy: %w", err)
	}
	if err := atomicfile.Write(path, starter, 0o644); err != nil {
		return "", fmt.Errorf("create working copy: %w", err)
	}
	return path, nil
}

// Reset overwrites the working copy for id with starter.
func (w Workspace) Reset(id string, starter []byte) (string, error) {
	if w.Dir == "" {
		return "", errNoDir
	}
	path := w.Path(id)
	if err := atomicfile.Write(path, starter, 0o644); err != nil {
		return "", fmt.Errorf("reset working copy: %w", err)
	}
	return path, nil
}
