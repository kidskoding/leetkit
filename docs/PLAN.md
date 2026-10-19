# LeetKit Implementation Plan

Source of truth: [`SPEC.md`](../SPEC.md). This plan records how we build it, in what order, and where we
deviate from the suggested structure.

## Verified dependencies (2026-10-19)

APIs were checked against the package docs (`go doc`) and the upgrade guides shipped in each module.

| Module | Version | Notes |
|---|---|---|
| `charm.land/bubbletea/v2` | v2.1.0 | `View() tea.View`, `tea.NewView`, `v.AltScreen = true`, `tea.KeyPressMsg`, `tea.ExecProcess(*exec.Cmd, func(error) tea.Msg)` |
| `charm.land/bubbles/v2` | v2.2.1 | `list.New(items, delegate, w, h)`, `viewport.New()` + `SetWidth/SetHeight`, `textinput.New()`, `help.New()`, `key.NewBinding`, `progress.New(...)`, `spinner.New(...)` |
| `charm.land/lipgloss/v2` | v2.0.6 | `lipgloss.Color("#hex")` returns `color.Color`; no `AdaptiveColor` |
| `github.com/BurntSushi/toml` | v1.6.0 | exercise metadata and config |

Bubble Tea v2 downsamples colors to the detected terminal profile, so fixed Tokyo Night hex colors degrade
gracefully on 256/16-color terminals.

## Architecture

```text
main.go                    entry point: calls cli.Run(os.Args[1:])
exercises/                 bundled curriculum, embedded with go:embed (exercises.FS)
  modules.toml             module ids, titles, order
  <module>/<id>/           exercise.toml, starter.py, test_exercise.py, solution.py, explanation.md
internal/
  config/                  config.toml, platform data/config dirs
  exercise/                types, loader (from fs.FS), validation
  workspace/               learner working copies under the data dir; reset
  runner/                  Python discovery, pytest subprocess, JUnit XML parsing -> Result
  editor/                  $VISUAL/$EDITOR resolution, safe argument splitting -> *exec.Cmd
  progress/                JSON store, atomic writes, corruption recovery
  cli/                     subcommands: list, run, check, progress, reset, doctor
  tui/                     Bubble Tea model, views, theme, keymap
```

Data flow for one exercise:

```text
embedded starter.py --(first open)--> <data>/workspace/<id>/starter.py   (learner edits this)
run: temp dir <- learner starter.py + packaged test_exercise.py
     python -m pytest --junitxml=report.xml  (exec.CommandContext, timeout)
     parse report.xml -> runner.Result -> progress.Record(id, passed)
```

### Key decisions

- **pytest output:** use built-in `--junitxml`. It is machine-readable and needs no plugin
  (`pytest-json-report` would add a dependency). Raw stdout/stderr are kept for traceback inspection.
- **Isolated test runs:** each run copies the learner file and the packaged tests into a fresh temp dir and
  sets `PYTHONDONTWRITEBYTECODE=1` and `-p no:cacheprovider`. The workspace never collects caches, and
  learners cannot accidentally change the tests.
- **Python discovery:** config `python` path, then the LeetKit venv (`<data>/venv`), then `python3`,
  then `python`. Nothing is installed automatically; `leetkit doctor` prints the venv commands.
- **Editor:** `$VISUAL`, then `$EDITOR`, then the first of `nvim`, `vim`, `vi`, `nano` on `PATH`
  (`notepad` on Windows). The value is split into argv with quote-aware splitting; no shell is used.
- **CLI parsing:** standard library `flag`, one `FlagSet` per subcommand. No Cobra.
- **No sandbox:** curriculum is trusted and bundled. Runs use a temp working dir and a timeout only.

### Deviations from SPEC §12 suggested layout

| Spec suggestion | What we do | Why |
|---|---|---|
| `cmd/leetkit/` and `main.go` | root `main.go` only | `go run .` and `go install <module>@latest` both produce `leetkit` |
| `internal/app` + `internal/ui` | single `internal/tui` | one Bubble Tea model; splitting adds indirection without benefit |
| — | `internal/workspace`, `internal/cli` | working-copy management and CLI need to be testable without a TUI |
| Interview relevance (unspecified location) | `interview` field in `exercise.toml` | lets validation require it |
| `exercises/patterns/` | not created in MVP | none of the 15 MVP exercises belong to module 7 |

## MVP exercise ids

| # | id | module | type |
|---|---|---|---|
| 01 | `enumerate-01` | iteration | syntax |
| 02 | `zip-01` | iteration | syntax |
| 03 | `dict-get-01` | hashing | application (frequency map) |
| 04 | `sets-01` | hashing | application (duplicates) |
| 05 | `counter-01` | collections | application (valid anagram) |
| 06 | `defaultdict-01` | collections | application (group anagrams) |
| 07 | `deque-01` | collections | application (BFS) |
| 08 | `list-comprehension-01` | pythonic | syntax |
| 09 | `dict-comprehension-01` | pythonic | syntax |
| 10 | `sorted-key-01` | sorting | syntax |
| 11 | `lambda-01` | sorting | application (intervals) |
| 12 | `tuple-unpacking-01` | iteration | syntax |
| 13 | `heapq-01` | utilities | application (top-k) |
| 14 | `bisect-01` | sorting | syntax |
| 15 | `accumulate-01` | utilities | application (prefix sums) |

## Milestones

Each milestone ends at a checkpoint: `go fmt ./... && go vet ./... && go test ./...` pass, then stop for a
commit.

0. **Scaffold** — plan, `CLAUDE.md`, `go.mod`, `.gitignore`.
1. **Exercise loading** — `exercise` package, `go:embed`, `modules.toml`, first exercise (`enumerate-01`),
   validation (missing fields, duplicate ids, bad module), `leetkit list`.
2. **Execution engine** — `config`, `workspace`, `runner` (discovery, pytest, JUnit parsing, timeout,
   cancellation, syntax errors), `leetkit check <id>`, `leetkit doctor`.
3. **Progress** — `progress` store (atomic write, corrupt-file recovery), record on check,
   `leetkit progress`, `leetkit reset <id> [--yes]`.
4. **TUI vertical slice** — list + detail panes, `e` opens editor via `tea.ExecProcess`, `r` runs tests in a
   `tea.Cmd`, results view, progress saved. **End-to-end acceptance point.**
5. **Learning features** — hints, explanation view, search (`/`), help (`?`), progress view (`p`),
   next (`n`), `Tab` focus, `leetkit run <id>`, small-terminal message.
6. **Curriculum** — remaining 14 exercises; Go test runs every reference solution (must pass) and every
   starter (must fail).
7. **Polish** — README, CONTRIBUTING, LICENSE, Makefile, release build config, error presentation.

## Known limitations (to revisit)

- Not a sandbox: exercise code runs with the user's permissions.
- Requires pytest; LeetKit will not install it.
