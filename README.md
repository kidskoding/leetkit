# LeetKit

**Master Python's toolkit for coding interviews.** Inspired by [Rustlings](https://github.com/rust-lang/rustlings).

LeetKit is a terminal trainer for the Python features that make interview solutions short and correct:
`enumerate`, `zip`, `Counter`, `defaultdict`, `deque`, `heapq`, `bisect`, comprehensions, slicing, and more.
As in Rustlings, you fix small exercises in your own editor while LeetKit watches the file, runs the tests,
and tracks your progress.

```text
 07/55  pairwise-01 · Largest Temperature Swing
 Iteration and Unpacking · beginner · Interview Application

 Return the largest absolute difference between two consecutive readings...

 Example 1
   Input:       largest_swing([12, 15, 9, 10])
   Output:      6
   Explanation: The consecutive differences are |15-12| = 3, |9-15| = 6 and
                |10-9| = 1. The largest is 6.

 ── Tests ──────────────────────────────────────────────────────────────
 ✗ 3 passed, 4 failed 240ms

   ✗ test_mixed_changes
       assert 3 == 6
        +  where 3 = largest_swing([12, 15, 9, 10])

   ✓ test_negative_readings
   ✓ test_flat_readings
   ...

 Progress ████████░░░░░░░░░░░░░░░░ 12/55
 Current  ~/.local/share/leetkit/workspace/pairwise-01/starter.py
 n/N next/prev • e edit • t hint 0/3 • r run • h/← list • ? more • q quit
```

## Why

Plenty of people can explain a sliding window or a BFS but still write `for i in range(len(nums))`, build
frequency maps by hand, or forget that `heapq` exists. Interviews reward fluent, idiomatic Python.
LeetKit drills exactly that: each exercise teaches one feature and then applies it to a small interview-style
problem. It is not an algorithms course and not a LeetCode replacement.

## Curriculum

55 exercises across seven modules, each with tests, progressive hints, worked examples, and an explanation
that unlocks once you solve it.

| Module | Covers |
|---|---|
| Iteration and Unpacking | `enumerate`, `zip`, `range`, `reversed`, tuple and star unpacking, `dict.items`, `zip(*matrix)` |
| Dictionaries and Sets | `dict.get`, set operators, `setdefault`, membership, one-pass Two Sum |
| Collections | `Counter`, `most_common`, `defaultdict`, `deque` for BFS and sliding windows |
| Sorting and Searching | `sorted(key=...)`, `lambda`, tuple keys, `min`/`max` with `key`, `bisect` |
| Pythonic Syntax | list/dict/set comprehensions, generator expressions, `any`/`all`, slicing, `str.join` |
| Interview Utilities | `heapq`, `math.inf`, `math.gcd`, `divmod`, `itertools.accumulate`, `itertools.combinations` |
| Applied Interview Patterns | two pointers, sliding window, prefix sums, BFS, DFS, binary search, heap selection |

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
4. When everything passes, press `x` for the explanation and `n` for the next exercise.

## Keys

| Key | Exercise screen | List |
|---|---|---|
| `j`/`k`, `↓`/`↑` | scroll | move |
| `l`, `→`, `Enter` | | open exercise |
| `h`, `←`, `Esc` | open the list | back |
| `n` / `N` | next / previous exercise | |
| `e` | edit in `$EDITOR` | |
| `r` | rerun tests | reset exercise (asks first) |
| `t` | reveal the next hint | |
| `x` | explanation (after solving) | |
| `/` | | search |
| `p` | progress by module | progress by module |
| `?` | all keys | all keys |
| `q` | quit | back |

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
