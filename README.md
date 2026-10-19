# LeetKit

[![CI](https://img.shields.io/github/actions/workflow/status/kidskoding/leetkit/ci.yml?branch=master&label=CI)](https://github.com/kidskoding/leetkit/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/kidskoding/leetkit?sort=semver)](https://github.com/kidskoding/leetkit/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/kidskoding/leetkit)](go.mod)
[![Python](https://img.shields.io/badge/python-3.10%2B-3776AB?logo=python&logoColor=white)](#install)
[![Exercises](https://img.shields.io/badge/exercises-91-2ea44f)](#curriculum)
[![Nix flake](https://img.shields.io/badge/nix-flake-5277C3?logo=nixos&logoColor=white)](flake.nix)
[![License](https://img.shields.io/github/license/kidskoding/leetkit)](LICENSE)

**Master Python's toolkit for coding interviews.** Inspired by [Rustlings](https://github.com/rust-lang/rustlings).

LeetKit is a terminal trainer for the Python features that make interview solutions short and correct:
`enumerate`, `zip`, `Counter`, `defaultdict`, `deque`, `heapq`, `bisect`, comprehensions, slicing, and more.
As in Rustlings, you fix small exercises in your own editor while LeetKit watches the file, runs the tests,
and tracks your progress.

```text
 Largest Temperature Swing                                               7 / 91
 pairwise-01 · Iteration and Unpacking · beginner · Interview Application
────────────────────────────────────────────────────────────────────────────────
 ── Tests ─────────────────────────────────────────────────────────────────────
 ✗ 3 passed, 4 failed 230ms

   ✗ test_mixed_changes
       assert 3 == 6
        +  where 3 = largest_swing([12, 15, 9, 10])

   ✓ test_negative_readings
   ✓ test_flat_readings
   ✓ test_only_adjacent_pairs_count
   ✗ test_fewer_than_two_readings
       ValueError: max() iterable argument is empty

────────────────────────────────────────────────────────────────────────────────
 ✗ 3 of 7 tests passing                           ██░░░░░░░░░░░░░░ 12/91 solved
 ~/.local/share/leetkit/workspace/pairwise-01/starter.py
 e edit • r run • t hint • n next • l list • ? help • q quit
```

## Why

Plenty of people can explain a sliding window or a BFS but still write `for i in range(len(nums))`, build
frequency maps by hand, or forget that `heapq` exists. Interviews reward fluent, idiomatic Python.
LeetKit drills exactly that: each exercise teaches one feature and then applies it to a small interview-style
problem. It is not an algorithms course and not a LeetCode replacement.

## Curriculum

91 exercises across seven modules, each with tests, progressive hints, worked examples, and an explanation
that unlocks once you solve it.

| Module | Covers |
|---|---|
| Iteration and Unpacking | `enumerate`, `zip` (and `strict=True`), `range`, `reversed`, tuple, star, and nested unpacking, `dict.items`, grid traversal |
| Dictionaries and Sets | `dict.get`, set operators, `setdefault`, `frozenset` keys, Two Sum, consecutive sequences, cycle detection |
| Collections | `Counter` (windows, `&`, `most_common`), `defaultdict`, `deque` (BFS, monotonic window), `OrderedDict` LRU cache, `namedtuple` |
| Sorting and Searching | `sorted(key=...)`, `lambda`, tuple keys, `functools.cmp_to_key`, `min`/`max` with `key`, `bisect`, `insort`, LIS |
| Pythonic Syntax | comprehensions, generator expressions, `any`/`all`, slicing, `str` methods, f-strings, `:=`, `*args`/`**kwargs`, `map`/`filter` |
| Interview Utilities | `heapq` (`merge`, `nlargest`), `math.inf`, `math.gcd`, `divmod`, `itertools` (`accumulate`, `combinations`, `product`, `permutations`, `groupby`), `functools.cache` |
| Applied Interview Patterns | two pointers, sliding window, prefix sums, BFS, DFS, binary search, heaps, monotonic stack, topological sort, backtracking, DP, union-find |

## Install

LeetKit needs **Python 3.10+ with pytest** to run exercises (the Nix package bundles both).

```bash
# Nix: no Python setup needed
nix run github:kidskoding/leetkit

# Go 1.26+
go install github.com/kidskoding/leetkit@latest
```

Prebuilt binaries for Linux, macOS, and Windows are attached to each
[release](https://github.com/kidskoding/leetkit/releases).

If you installed with Go or a release binary, give LeetKit a Python with pytest. `leetkit doctor` checks your
setup and prints these commands for your machine:

```bash
uv venv ~/.local/share/leetkit/venv
uv pip install --python ~/.local/share/leetkit/venv pytest
# or, without uv:
python3 -m venv ~/.local/share/leetkit/venv
~/.local/share/leetkit/venv/bin/python -m pip install pytest
```

LeetKit never installs anything itself.

## Quick start

```bash
leetkit doctor   # check Python, pytest, and your editor
leetkit          # open the dashboard on your current exercise
```

1. Read the problem and examples.
2. Press `e` to open the exercise in `$EDITOR`, or edit the file shown next to **Current** in another pane.
3. Save. LeetKit reruns the tests automatically.
4. The status bar shows the result at all times. When everything passes, press `x` for the explanation and
   `n` for the next exercise. `l` opens the full list.

## Keys

On the exercise screen:

| Key | |
|---|---|
| `e` | edit in `$EDITOR` (tests rerun when you save or exit) |
| `r` | rerun the tests |
| `t` | reveal the next hint |
| `n` / `N` | next / previous exercise |
| `x` | explanation, once solved |
| `l` | exercise list |
| `p` | progress by module |
| `j`/`k`, `↓`/`↑` | scroll |
| `?` | all keys |
| `q` | quit |

In the list: `j`/`k` move, `enter` or `l` opens, `/` searches, `r` resets an exercise (asks first), and
`esc`, `h`, or `q` go back.

## Commands

| Command | |
|---|---|
| `leetkit` | open the dashboard |
| `leetkit run <id>` | open the dashboard on one exercise |
| `leetkit list [--ids]` | list modules and exercises |
| `leetkit check [-v] [<id>]` | run one exercise's tests, or check every exercise you have started |
| `leetkit progress` | completion by module |
| `leetkit reset [-y] <id>` | restore the starter code and clear progress |
| `leetkit doctor` | check Python, pytest, editor, and data directory |
| `leetkit completion bash\|zsh\|fish` | print a shell completion script |
| `leetkit version` | print the version |

Shell completion, for example in fish:

```bash
leetkit completion fish > ~/.config/fish/completions/leetkit.fish
```

## Configuration

Optional, in `~/.config/leetkit/config.toml` (`%AppData%\leetkit` on Windows,
`~/Library/Application Support/leetkit` on macOS):

```toml
python  = "/path/to/python3"   # default: $LEETKIT_PYTHON, LeetKit's venv, then python3 on PATH
editor  = "code --wait"        # default: $VISUAL, then $EDITOR, then nvim/vim/vi/nano
timeout = "10s"                # per test run
theme   = "terminal"           # terminal (follows your terminal colors), gruvbox, gruvbox-light, tokyonight
```

Your code and progress live in `~/.local/share/leetkit` on Linux (`$XDG_DATA_HOME` is honoured),
`~/Library/Application Support/leetkit` on macOS, and `%LocalAppData%\leetkit` on Windows. The bundled
exercises are never modified.

## Exercise structure

```text
exercises/
  01_iteration/
    01_enumerate-01/
      exercise.toml      metadata, description, hints, examples
      starter.py         what you start from
      test_exercise.py   pytest tests
      solution.py        reference solution
      explanation.md     the feature and its interview relevance
```

Directory prefixes set the order. See [CONTRIBUTING.md](CONTRIBUTING.md) to write exercises.

## Security

Exercise code runs as a normal Python process with your permissions. LeetKit uses a scratch directory and a
timeout, but it is **not a sandbox**. The bundled curriculum is trusted; review any third-party exercises
before running them.

## Development

```bash
direnv allow     # Nix devShell: Go, Python + pytest, uv  (or: uv sync)
go test ./...    # includes running every exercise's solution and starter
go run .
```

Details in [CONTRIBUTING.md](CONTRIBUTING.md).

## Roadmap

- More exercises, especially interview applications that combine several tools
- Optional reference solution after repeated failed attempts
- Third-party exercise packs

## Acknowledgements

LeetKit is inspired by [Rustlings](https://github.com/rust-lang/rustlings): the watch-mode workflow, the
numbered exercise layout, and the idea of learning a language through small failing exercises all come from
it. [Pylings](https://github.com/CompEng0001/Pylings) brought the same approach to Python. LeetKit is built
with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles),
and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## License

[MIT](LICENSE). LeetKit is an independent project and is not affiliated with or endorsed by LeetCode.
