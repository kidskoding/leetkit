// Package runner executes exercise tests with pytest in a subprocess.
//
// Learner code runs with the user's permissions. Runs use a scratch directory
// and a timeout, but this is not a sandbox.
package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
)

const (
	reportFile     = "report.xml"
	maxOutputBytes = 1 << 20
	// pytest exit code when no tests were collected.
	exitNoTests = 5
)

type Runner struct {
	Python  string
	Timeout time.Duration
}

type Job struct {
	Tests []byte
	Code  []byte
	// CodePath is the learner's file; it replaces scratch paths in messages.
	CodePath string
}

// Run returns a Result for anything the learner's code does, including
// timeouts. It returns an error only when tests could not be started or
// pytest itself is unusable.
func (r Runner) Run(ctx context.Context, job Job) (*Result, error) {
	dir, err := os.MkdirTemp("", "leetkit-run-")
	if err != nil {
		return nil, fmt.Errorf("create scratch directory: %w", err)
	}
	defer os.RemoveAll(dir)

	files := map[string][]byte{
		exercise.StarterFile: job.Code,
		exercise.TestFile:    job.Tests,
		// An ini file pins pytest's rootdir here so unrelated parent configs are ignored.
		"pytest.ini": []byte("[pytest]\n"),
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			return nil, fmt.Errorf("prepare test run: %w", err)
		}
	}

	if r.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	output := &cappedBuffer{limit: maxOutputBytes}
	cmd := exec.CommandContext(ctx, r.Python, "-B", "-m", "pytest",
		"-p", "no:cacheprovider", "-q", "--tb=short", "--color=no",
		"--junitxml="+reportFile, exercise.TestFile)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PYTEST_DISABLE_PLUGIN_AUTOLOAD=1", "PYTEST_ADDOPTS=")
	cmd.Stdout, cmd.Stderr = output, output
	cmd.WaitDelay = time.Second

	start := time.Now()
	runErr := cmd.Run()
	res := &Result{Duration: time.Since(start)}
	clean := strings.NewReplacer(
		filepath.Join(dir, exercise.StarterFile), job.CodePath,
		dir+string(filepath.Separator), "",
	)
	res.Output = clean.Replace(output.String())

	switch {
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		res.Status = StatusTimedOut
		res.Problem = fmt.Sprintf("Tests did not finish within %s. Look for an infinite loop.", r.Timeout)
		return res, nil
	case errors.Is(ctx.Err(), context.Canceled):
		res.Status = StatusCanceled
		res.Problem = "Test run was canceled."
		return res, nil
	}

	var exitErr *exec.ExitError
	if runErr != nil && !errors.As(runErr, &exitErr) {
		return nil, fmt.Errorf("start %s: %w", r.Python, runErr)
	}

	report, err := os.ReadFile(filepath.Join(dir, reportFile))
	if err != nil {
		if strings.Contains(res.Output, "No module named pytest") {
			return nil, ErrPytestMissing
		}
		return nil, fmt.Errorf("pytest exited without a report (%v):\n%s", runErr, strings.TrimSpace(res.Output))
	}
	if res.Cases, err = parseJUnit([]byte(clean.Replace(string(report)))); err != nil {
		return nil, err
	}

	switch {
	case res.Count(OutcomeError) > 0:
		res.Status = StatusError
		res.Problem = "Your code could not be tested. Fix the error below and run again."
	case exitErr != nil && exitErr.ExitCode() == exitNoTests, len(res.Cases) == 0:
		res.Status = StatusError
		res.Problem = "No tests were collected. The exercise's test file may be broken."
	case res.Count(OutcomeFailed) > 0 || runErr != nil:
		res.Status = StatusFailed
	default:
		res.Status = StatusPassed
	}
	return res, nil
}

// cappedBuffer keeps the first limit bytes so runaway output cannot exhaust memory.
type cappedBuffer struct {
	buf       bytes.Buffer
	limit     int
	truncated bool
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	if room := b.limit - b.buf.Len(); room < len(p) {
		b.truncated = true
		b.buf.Write(p[:max(room, 0)])
		return len(p), nil
	}
	return b.buf.Write(p)
}

func (b *cappedBuffer) String() string {
	if b.truncated {
		return b.buf.String() + "\n[output truncated]\n"
	}
	return b.buf.String()
}
