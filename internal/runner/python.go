package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	ErrNoPython      = errors.New("no Python interpreter found")
	ErrPytestMissing = errors.New("pytest is not installed for this interpreter")
)

// Candidate is a possible interpreter. Explicit candidates come from user
// settings: if one is set but unusable, discovery fails instead of silently
// falling back to a different Python.
type Candidate struct {
	Path     string
	Source   string
	Explicit bool
}

// FindPython returns the first usable candidate, then falls back to python3
// and python on PATH. Candidates with an empty Path are ignored.
func FindPython(candidates ...Candidate) (Candidate, error) {
	for _, c := range candidates {
		if c.Path == "" {
			continue
		}
		if path, err := exec.LookPath(c.Path); err == nil {
			c.Path = path
			return c, nil
		}
		if c.Explicit {
			return Candidate{}, fmt.Errorf("%s is set to %q, which is not an executable file", c.Source, c.Path)
		}
	}
	for _, name := range []string{"python3", "python"} {
		if path, err := exec.LookPath(name); err == nil {
			return Candidate{Path: path, Source: "PATH"}, nil
		}
	}
	return Candidate{}, ErrNoPython
}

func PythonVersion(ctx context.Context, python string) (string, error) {
	out, err := pythonOutput(ctx, python, "import platform; print(platform.python_version())")
	if err != nil {
		return "", fmt.Errorf("run %s: %w", python, err)
	}
	return out, nil
}

func PytestVersion(ctx context.Context, python string) (string, error) {
	out, err := pythonOutput(ctx, python, "import pytest; print(pytest.__version__)")
	if err != nil {
		if strings.Contains(err.Error(), "No module named 'pytest'") {
			return "", ErrPytestMissing
		}
		return "", fmt.Errorf("check pytest: %w", err)
	}
	return out, nil
}

func pythonOutput(ctx context.Context, python, code string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, python, "-c", code)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		if msg := lastLine(stderr.String()); msg != "" {
			return "", fmt.Errorf("%w: %s", err, msg)
		}
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func lastLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.LastIndexByte(s, '\n'); i >= 0 {
		return s[i+1:]
	}
	return s
}
