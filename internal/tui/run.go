package tui

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
	"github.com/anirudhkonidala/leetkit/internal/runner"
	"github.com/anirudhkonidala/leetkit/internal/workspace"
)

type runState struct {
	active   bool
	seq      int
	cancelFn context.CancelFunc
	result   *runner.Result
	err      error
}

func (r runState) cancel() {
	if r.cancelFn != nil {
		r.cancelFn()
	}
}

type requestRunMsg struct{}

// runDoneMsg carries the sequence number of the run it answers, so results
// from a superseded run are ignored.
type runDoneMsg struct {
	seq int
	id  string
	res *runner.Result
	err error
}

// startRun cancels any run in progress and tests the current exercise in the
// background. Only one run exists at a time.
func (m *Model) startRun() tea.Cmd {
	ex := m.currentExercise()
	if ex == nil {
		return nil
	}
	m.run.cancel()
	m.run = runState{seq: m.run.seq + 1}
	if m.opts.RunnerErr != nil {
		m.run.err = m.opts.RunnerErr
		m.refreshWatch()
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.run.active, m.run.cancelFn = true, cancel
	m.refreshWatch()

	seq, cat, ws, r := m.run.seq, m.opts.Catalog, m.opts.Workspace, m.opts.Runner
	run := func() tea.Msg {
		job, err := prepareJob(cat, ws, ex)
		if err != nil {
			return runDoneMsg{seq: seq, id: ex.ID, err: err}
		}
		res, err := r.Run(ctx, job)
		return runDoneMsg{seq: seq, id: ex.ID, res: res, err: err}
	}
	return tea.Batch(m.spinner.Tick, run)
}

func prepareJob(cat *exercise.Catalog, ws workspace.Workspace, ex *exercise.Exercise) (runner.Job, error) {
	starter, err := cat.ReadFile(ex, exercise.StarterFile)
	if err != nil {
		return runner.Job{}, fmt.Errorf("read starter: %w", err)
	}
	tests, err := cat.ReadFile(ex, exercise.TestFile)
	if err != nil {
		return runner.Job{}, fmt.Errorf("read tests: %w", err)
	}
	path, err := ws.Ensure(ex.ID, starter)
	if err != nil {
		return runner.Job{}, err
	}
	code, err := os.ReadFile(path)
	if err != nil {
		return runner.Job{}, fmt.Errorf("read your code: %w", err)
	}
	return runner.Job{Tests: tests, Code: code, CodePath: path}, nil
}

func (m *Model) finishRun(msg runDoneMsg) tea.Cmd {
	if msg.seq != m.run.seq {
		return nil
	}
	m.run.cancel()
	m.run.active, m.run.result, m.run.err = false, msg.res, msg.err
	if msg.err == nil && msg.res.Status != runner.StatusCanceled {
		if _, err := m.opts.Store.Record(msg.id, msg.res.Status == runner.StatusPassed); err != nil {
			m.notice = "Could not save progress: " + err.Error()
		}
		m.refreshItems()
	}
	m.refreshWatch()
	return nil
}
