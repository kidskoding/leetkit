# LeetKit

Terminal app (Go + Charm v2) that trains Python syntax and standard-library fluency for coding interviews.
Learners edit Python exercises in their own editor; LeetKit runs pytest and tracks progress.

- **Spec (source of truth):** `SPEC.md`
- **Deviations from the spec:** see "Deviations from SPEC.md" below; update it when a decision changes.

## Workflow rules

- Work in small units (one package or one user-facing feature). Commit each unit with a Conventional Commits
  message, no `Co-Authored-By` trailer, then stop for review. Every commit must pass `go vet` and `go test` on
  its own. **Never `git push`.**
- Smallest working vertical slice first. No feature is done until it works end to end.
- Write tests with the code. Business logic must be testable without a terminal.
- Record any departure from `SPEC.md` under "Deviations from SPEC.md" below.

## Commands

```bash
go run .                       # launch the TUI
go run . list [--ids] | run <id> | check [-v] [<id>] | progress | reset [-y] <id> | doctor | completion <shell>
go fmt ./... && go vet ./... && go test ./...   # required before every checkpoint
go test -race ./...
```

## Python environment

LeetKit needs one thing: a Python 3.10+ interpreter with pytest. It never installs anything itself.

- **Nix:** `flake.nix` + `default.nix` (gomod2nix). `nix run` / `nix build` produce a binary with a bundled
  Python + pytest (linked in via `-X main.defaultPython`). direnv (`.envrc`) loads the devShell (Go, Python +
  pytest, uv). After changing Go dependencies, run `gomod2nix generate` (the devShell does it automatically)
  and commit `gomod2nix.toml`.
- **uv (contributors, CI):** `uv sync` creates `./.venv` from `pyproject.toml` / `uv.lock`.
- **End users:** `leetkit doctor` prints venv setup commands for `<data>/venv`.

Interpreter lookup order: `$LEETKIT_PYTHON`, config `python`, `<data>/venv`, packaged default, `python3`,
`python`.
An explicitly set interpreter that is missing is an error, never a silent fallback.
Go tests that need pytest prefer `./.venv/bin/python`, then `PATH`, and skip if neither has pytest.

## Layout

```text
main.go              -> cli.Run(os.Args[1:])
exercises/           bundled curriculum, embedded via go:embed (exercises.FS); modules.toml lists modules
internal/config      config.toml + platform data/config dirs
internal/exercise    types, loader from fs.FS, validation
internal/workspace   learner working copies in <data>/workspace/<id>/
internal/runner      Python discovery, pytest subprocess, JUnit XML -> Result
internal/editor      $VISUAL/$EDITOR resolution -> *exec.Cmd
internal/progress    JSON progress store
internal/cli         subcommands (stdlib flag)
internal/theme       named palettes (role-based) and the styles built from them
internal/tui         Bubble Tea model, views, keys
```

Dependency direction: `tui`/`cli` -> `runner`, `workspace`, `progress`, `exercise`, `editor`, `config`.
Lower packages never import `tui` or `cli`.

## Charm v2 rules (verified against v2 docs — do not use v1 patterns)

- Imports: `tea "charm.land/bubbletea/v2"`, `"charm.land/lipgloss/v2"`, `"charm.land/bubbles/v2/<pkg>"`.
  Never `github.com/charmbracelet/bubbletea` etc.
- `View() tea.View`, built with `tea.NewView(s)`. Set `v.AltScreen = true` on the view, not a program option.
- Keys: match `tea.KeyPressMsg` and `msg.String()` (`"ctrl+c"`, `"enter"`, `"space"`). `tea.KeyMsg` is an
  interface that also matches key releases.
- Bubbles components use setters: `vp.SetWidth(w)`, `vp.SetHeight(h)`, `help.SetWidth`, `textinput.SetWidth`.
  `viewport.New()` takes options, not `(w, h)`. `list.New(items, delegate, w, h)`.
- `AdaptiveColor` no longer exists; use `lipgloss.LightDark(isDark)` if a style must vary by background.
- Never hardcode colors in views. Take a `theme.Theme` and use its styles. The default `terminal` theme uses
  ANSI palette indices so it matches the user's terminal scheme; hex themes are opt-in via `theme` in config.
- CLI output goes through `colorprofile.NewWriter` (set up in `main.go`), which downsamples and strips color
  for pipes and `NO_COLOR`. Tests use the same writer, so assertions see plain text.
- Editor: `tea.ExecProcess(cmd, func(err error) tea.Msg { ... })`. Slow work (pytest) runs in a `tea.Cmd`
  that returns a typed message. Never block `Update`.
- Key handling must respect focus: when the search input is focused, letters go to the input, not shortcuts.

## Deviations from SPEC.md

- Default theme is `terminal` (ANSI colors follow the user's terminal) instead of Tokyo Night; `tokyonight`,
  `gruvbox`, `gruvbox-light` are selectable with `theme = "..."` in config.toml.
- `leetkit check` without an id checks every started exercise.
- The dashboard follows Rustlings: it opens on the current exercise (watch screen) and the exercise list is a
  secondary screen. The list is a flat `bubbles/list` (module shown in the preview) rather than a collapsible
  module tree, which keeps Bubbles' built-in fuzzy search.
- Keys use ranger/lf-style navigation: j/k move, l/→/enter open an exercise, h/←/esc go to the list. Hints are
  on `t` (SPEC: `h`) and reset is `r` in the list. `n`/`N` cycle forward/backward through all exercises
  (wrapping), solved or not; Rustlings only allows `n` after solving.
- The curriculum covers all seven modules (91 exercises and growing) instead of SPEC's 15-exercise MVP, at the
  maintainer's request. `go test ./exercises` verifies every solution passes and every starter fails.
- Added beyond SPEC: tests rerun automatically when the working copy is saved (mtime polling, like Rustlings
  watch mode).
- Root `main.go` instead of `cmd/leetkit`; single `internal/tui` instead of `internal/app` + `internal/ui`.

## Go conventions

- Standard library first. Current external deps: Charm v2 trio, `github.com/BurntSushi/toml`. Ask before
  adding others.
- Wrap errors with context (`fmt.Errorf("load %s: %w", path, err)`); never swallow them. CLI errors must say
  how to fix the problem.
- Subprocesses: `exec.CommandContext` with explicit argv. Never `sh -c` or string concatenation.
- File writes that matter (progress, config) are atomic: write temp file in the same dir, then `os.Rename`.
- Comments only where they add information: package docs, the why behind a non-obvious choice, or behavior
  a caller can't infer from the signature. No doc comments that restate the identifier name.
- No speculative interfaces. Introduce an interface only when a second implementation or a test needs it.

## Exercise authoring

Each exercise lives in `exercises/NN_<module>/NN_<id>/`. The `NN_` prefixes set the order of modules and of
exercises within a module (Rustlings-style); there is no `order` field. `modules.toml` maps module ids to
titles. Renumbering never affects learners: workspaces and progress are keyed by `id`.

| File | Purpose |
|---|---|
| `exercise.toml` | `id` (directory name without `NN_`), `title`, `module`, `difficulty`, `type` (`syntax`/`application`), `concepts`, `description`, `objective`, `interview`, `hints`, `[[examples]]` with `input`, `output`, `explanation` |
| `starter.py` | function signature, one `# TODO:` comment, `pass`; must import cleanly and fail its tests |
| `test_exercise.py` | 4-8 pytest tests importing from `starter`; normal and edge cases |
| `solution.py` | idiomatic reference solution; must pass the tests |
| `explanation.md` | `#` title, what the feature does, `## Useful variations`, `## Interview relevance` |

- The packaged `starter.py` is a template. Learners edit a copy in the workspace, which LeetKit prefixes with
  a comment header generated from `exercise.toml` (title, description, examples). Never write to
  `exercises/`.
- Examples are LeetCode-style: `input` is a call expression, `output` the repr of the result, `explanation`
  one or two sentences on why. Prefer two examples, the second showing an edge case.
- Tests check correctness first. Syntax drills may add one feature-usage check with `ast` (for example,
  `ast.ListComp` or a call to `enumerate`) whose assertion message explains the requirement. No
  source-string grepping.
- `explanation.md` uses only headings, paragraphs, fenced code, `- ` bullets, and inline code; the dashboard
  renders that subset.
- `go test ./exercises` runs every reference solution (must pass) and every starter (must fail cleanly).
  Run one exercise with `go test ./exercises -run 'TestCurriculum/<id>'`.

## Safety

LeetKit runs Python with the user's permissions. It is not a sandbox; do not claim it is. Remote exercise
downloads are out of scope.
