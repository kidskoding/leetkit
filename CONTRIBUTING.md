# Contributing to LeetKit

Thanks for helping. The most valuable contributions are new exercises and fixes to existing ones; the Go code
is there to serve them.

## Development setup

You need Go 1.26+ and Python 3.10+ with pytest. Pick one:

```bash
# Nix + direnv: Go, Python 3.15 with pytest, uv, and Go tooling
direnv allow            # or: nix develop

# uv
uv sync                 # creates .venv with pytest; tests find it automatically
```

Then:

```bash
go run .                # the dashboard
go test ./...           # everything, including every exercise
go test -race ./...
gofmt -l . && go vet ./...
```

Tests that need pytest look for `./.venv/bin/python`, then `python3` on `PATH`, and skip when neither has
pytest. CI runs them all on Linux and macOS.

## Writing an exercise

Each exercise teaches **one Python feature** and applies it to a small, realistic interview-style problem.
It should take 2-5 minutes. Read a few existing exercises first, for example
`exercises/01_iteration/01_enumerate-01/`.

### Where it goes

```text
exercises/NN_<module>/NN_<id>/
  exercise.toml
  starter.py
  test_exercise.py
  solution.py
  explanation.md
```

- `NN_` prefixes set the order. To add an exercise at the end of a module, use the next free number. To
  insert one, renumber the directories after it; ids, saved work, and progress are unaffected.
- The id is the directory name without its prefix, e.g. `enumerate-01`. Ids are unique across the
  curriculum. Use `<topic>-NN`.
- New modules need a directory `NN_<id>/` and an entry in `exercises/modules.toml`.

### exercise.toml

```toml
id = "zip-01"
title = "Grade an Answer Sheet"
module = "iteration"
difficulty = "beginner"     # beginner | intermediate | advanced
type = "syntax"             # syntax: drill one feature | application: small interview problem
concepts = ["zip", "iteration"]

description = """
Implement count_correct(answers, key).

Return how many positions hold the same value in both lists. If the lists
have different lengths, compare only up to the shorter one.

Use zip() instead of indexing.
"""

objective = "Walk two sequences in lockstep with zip()."

interview = """
Two to four sentences on where this shows up in interviews and why the idiom matters.
"""

hints = [
  "A gentle nudge toward the feature.",
  "Something more concrete: name the API or pattern.",
  "Nearly the answer, without giving full code.",
]

[[examples]]
input = "count_correct(['a', 'c', 'b', 'd'], ['a', 'b', 'b', 'd'])"
output = "3"
explanation = "Positions 0, 2 and 3 match the key ('a', 'b', 'd'); position 1 has 'c' instead of 'b'. That is 3 correct."
```

- The description says exactly what to return, including edge cases (empty input, ties, duplicates,
  ordering), so the tests are fair.
- Examples are LeetCode-style: `input` is a call expression, `output` is the `repr` of the result, and
  `explanation` is one or two sentences on why. Prefer two examples, the second showing an edge case.
  Avoid outputs whose `repr` order can vary, such as sets with several elements.
- Unknown keys are rejected, so typos fail loudly.

### starter.py

The function signature, one `# TODO:` comment saying what to do, and `pass`. Include imports the intended
solution needs (`from collections import Counter`) so the learner sees the API. It must import cleanly and
fail the tests. No other comments or docstrings: LeetKit adds the problem statement as a header when it
creates the learner's copy.

### test_exercise.py

- `from starter import <names>`, plain pytest `assert`s, no I/O, no randomness.
- 4-8 tests named for what they check (`test_empty_list`, `test_ties_keep_input_order`). Cover normal cases
  and the edge cases the description promises.
- Correctness first. A **syntax drill** may add one feature-usage check using `ast`, never by searching the
  source text (comments must not count). Explain the requirement in the assertion message:

  ```python
  import ast
  import inspect

  def uses_name(func, name):
      tree = ast.parse(inspect.getsource(func))
      return any(
          (isinstance(node, ast.Name) and node.id == name)
          or (isinstance(node, ast.Attribute) and node.attr == name)
          for node in ast.walk(tree)
      )

  def test_uses_zip():
      assert uses_name(count_correct, "zip"), "This drill practices zip(): loop with `for a, b in zip(xs, ys)`."
  ```

  For comprehension drills check node types instead (`ast.ListComp`, `ast.DictComp`, `ast.SetComp`,
  `ast.GeneratorExp`).

### solution.py

Short, idiomatic, and uses the feature being taught. No comments.

### explanation.md

Shown after the learner solves the exercise:

```markdown
# zip()

What it does, with a short code block, and the clunky version it replaces.

## Useful variations

- Two to four bullets with inline `code`.

## Interview relevance

Where it appears in interviews and why it matters.
```

The dashboard renders headings, paragraphs, fenced code, `- ` bullets, and inline code. Avoid tables, links,
and images.

### Check it

```bash
go test ./exercises -run 'TestCurriculum/zip-01'   # one exercise
go test ./exercises                                 # all of them
```

This verifies the metadata, that `solution.py` passes, and that `starter.py` fails with ordinary test
failures rather than a syntax or import error. Then try it as a learner:

```bash
go run . run zip-01
```

## Code changes

- Read `CLAUDE.md` for architecture, conventions, and the Charm v2 rules (Bubble Tea, Bubbles, and Lip Gloss
  v2 APIs only).
- Standard library first; discuss new dependencies in an issue before adding them. After changing Go
  dependencies, run `gomod2nix generate` (the devShell does it automatically) and commit `gomod2nix.toml`.
- Write tests with the change. The TUI model is tested without a terminal by sending messages to `Update`.
- Keep commits small and focused, with [Conventional Commits](https://www.conventionalcommits.org/) messages
  (`feat(tui): ...`, `fix(runner): ...`, `docs: ...`). Every commit should pass `go vet` and `go test`.

## Releases

Maintainers tag `vX.Y.Z`; the release workflow builds binaries with GoReleaser and publishes them on GitHub.
