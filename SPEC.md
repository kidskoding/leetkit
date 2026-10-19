# LeetKit — Claude Code Project Specification

## 1. Project Overview

**Name:** LeetKit  
**CLI command:** `leetkit`  
**Tagline:** Master Python's toolkit for coding interviews.

LeetKit is an open-source, interactive terminal application designed to help developers master Python syntax, built-in functions, standard-library utilities, and idiomatic programming patterns frequently used in technical coding interviews.

The project is inspired by:

- [Rustlings](https://github.com/rust-lang/rustlings) — hands-on programming exercises
- [Pylings](https://github.com/CompEng0001/Pylings) — interactive Python learning
- LeetCode-style interview preparation — practical applications of Python features

Unlike general-purpose Python tutorials, LeetKit focuses exclusively on **Python fluency for coding interviews**.

### The Problem

Many developers understand algorithms and data structures but struggle to use Python effectively during interviews.

For example, they might know how to:

- Count character frequencies, but not know `Counter`.
- Group values in a dictionary, but not know `defaultdict`.
- Traverse graphs, but not know `deque`.
- Implement binary search manually, but not know `bisect`.
- Iterate through arrays using indices, but not know `enumerate`.
- Implement priority queues, but not know `heapq`.

LeetKit bridges the gap between understanding an algorithm and implementing it efficiently in Python.

### Target Audience

Developers who already understand basic programming and want to become fluent in Python for coding interviews.

This is not intended to be a beginner Python course or a full LeetCode replacement.

### Core Philosophy

**Learn the tool. Apply the tool. Understand when to use it.**

Every exercise should teach a practical Python feature and reinforce it with realistic coding interview scenarios.

---

## 2. Technology Stack

The application must be implemented in Go using the Charm terminal UI ecosystem.

| Component | Technology |
|---|---|
| Core language | Go |
| TUI framework | Bubble Tea v2 |
| UI styling | Lip Gloss v2 |
| Interactive widgets | Bubbles v2 |
| Exercise runtime | Python 3 |
| Python testing | pytest |
| Python execution | Go `os/exec` |
| Exercise metadata | TOML |
| Local persistence | JSON |
| Configuration | TOML |
| Code editor | External `$VISUAL` / `$EDITOR` |
| Default color scheme | Tokyo Night |

Use the official Charm v2 import paths:

```go
import (
    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
)
```

Use `charm.land/bubbles/v2` for UI components.

Follow the current Bubble Tea v2 API, including `View() tea.View` and its declarative terminal configuration.

Do not mix Bubble Tea v1 and v2 patterns.

Use current compatible stable dependency versions.

### Python Integration

Do not embed Python inside the Go application.

Do not use PyO3, cgo, or foreign-function interfaces.

Instead, launch a configurable Python interpreter as an external subprocess using `os/exec`.

The application is responsible for:

1. Discovering exercises.
2. Managing learner working files.
3. Opening an external editor.
4. Executing Python tests.
5. Capturing test results.
6. Displaying feedback.
7. Recording progress.

Python handles the actual learning exercises and their verification.

---

## 3. Product Experience

LeetKit should feel like a polished developer tool.

The interface must be:

- Fully keyboard-driven.
- Responsive to terminal resizing.
- Visually consistent.
- Fast to navigate.
- Friendly to Neovim users.
- Offline-first.
- Easy to extend.

### Main Dashboard

Launching:

```bash
leetkit
```

opens an interactive TUI.

Suggested layout:

```text
┌──────────────────────────────────────────────────────────────┐
│  LEETKIT                            Python Interview Toolkit │
│  Progress: ███████░░░░░░░░░░░  12/45                         │
├──────────────────────┬───────────────────────────────────────┤
│ MODULES / EXERCISES  │ EXERCISE DETAILS                      │
│                      │                                       │
│ ▼ Iteration          │ 03 — enumerate()                      │
│   ✓ Basics           │                                       │
│   ✓ zip()            │ Learn to iterate over indices and     │
│   > enumerate()      │ values simultaneously.                │
│   ○ reversed()       │                                       │
│                      │ Difficulty: Beginner                  │
│ ▶ Hash Maps          │ Type: Syntax Drill                    │
│ ▶ Collections        │                                       │
│ ▶ Sorting            │ [e] Edit exercise                     │
│ ▶ Pythonic Syntax    │ [r] Run tests                         │
│ ▶ Utilities          │ [h] Reveal hint                       │
│ ▶ Interview Patterns │                                       │
├──────────────────────┴───────────────────────────────────────┤
│ ↑↓ Navigate  Enter Select  e Edit  r Run  ? Help  q Quit     │
└──────────────────────────────────────────────────────────────┘
```

Treat this as a conceptual layout, not an exact rendering specification.

### Exercise Detail View

When an exercise is selected, display:

- Title
- Difficulty
- Exercise type
- Learning objective
- Problem description
- Example input and output
- Starter code preview
- Available hints
- Completion status

Use scrollable viewports for longer content.

### Test Results View

After running an exercise, present:

- Number of passed tests
- Number of failed tests
- Relevant assertion errors
- Python syntax errors
- Exception messages
- Execution time
- Overall pass/fail status

Avoid dumping raw pytest output without formatting.

Preserve the underlying traceback details so users can inspect them when necessary.

### Success State

When all tests pass:

```text
✓ Exercise Complete!

All 5 tests passed.

You've learned how to use enumerate().

[e] Edit again
[x] View explanation
[n] Next exercise
```

Keep animations subtle and optional.

### Keyboard Shortcuts

| Key | Behavior |
|---|---|
| `j` / `↓` | Navigate down |
| `k` / `↑` | Navigate up |
| `Enter` | Select |
| `Esc` | Return |
| `Tab` | Switch panel |
| `e` | Edit current exercise |
| `r` | Run tests |
| `h` | Reveal next hint |
| `x` | View explanation |
| `n` | Next exercise |
| `p` | Progress |
| `/` | Search |
| `?` | Help |
| `q` | Quit from primary navigation |

Shortcuts must respect focus. For example, typing `q` in a search field must not exit the application.

---

## 4. Learning System

The curriculum should teach specific Python language features rather than unrelated algorithms.

### Exercise Types

Every exercise belongs to one of two categories.

**A. Syntax Drill**

A short exercise focused on learning an individual Python feature.

Example: Use `enumerate()` to locate a target.

**B. Interview Application**

A small interview-style problem that applies the feature.

Example: Use `Counter` to determine whether two strings are anagrams.

Both exercise types should be available in the same curriculum.

Each exercise should generally take 2–5 minutes, although interview applications may take longer.

### Exercise Progression

Exercises should follow this pattern:

1. Introduce a Python feature.
2. Present a small incomplete implementation.
3. Ask the learner to complete the code.
4. Run automated tests.
5. Provide targeted feedback.
6. Explain the feature and its interview relevance.
7. Reinforce the feature through a related application.

Do not require every user to progress strictly sequentially.

---

## 5. Curriculum

Organize the curriculum into seven modules.

### Module 1: Iteration and Unpacking

Teach:

- `enumerate()`
- `zip()`
- `range()`
- `reversed()`
- Tuple unpacking
- Multiple assignment
- Iterating over dictionary items

Applications:

- Index tracking
- Comparing arrays
- Pairwise traversal
- Swapping variables
- Matrix iteration

### Module 2: Dictionaries and Sets

Teach:

- Dictionary creation and updates
- `dict.get()`
- `dict.items()`
- `dict.setdefault()`
- Set membership
- Set operations
- Dictionary comprehensions

Applications:

- Two Sum
- Duplicate detection
- Frequency maps
- Grouping and indexing

### Module 3: Collections

Teach:

- `collections.Counter`
- `collections.defaultdict`
- `collections.deque`
- `Counter.most_common()`
- Frequency counting
- Grouping values

Applications:

- Valid Anagram
- Group Anagrams
- BFS queue management
- Frequency-based problems

### Module 4: Sorting and Searching

Teach:

- `sorted()`
- `list.sort()`
- `lambda`
- Custom sorting keys
- `min()` and `max()` with `key`
- `bisect_left()`
- `bisect_right()`

Applications:

- Interval ordering
- Sorting coordinate pairs
- Binary-search insertion points
- Custom comparisons

### Module 5: Pythonic Syntax

Teach:

- List comprehensions
- Dictionary comprehensions
- Set comprehensions
- Generator expressions
- `any()`
- `all()`
- List slicing
- String joining
- Conditional expressions

Applications:

- Filtering
- Matrix construction
- Sequence transformation
- String manipulation

### Module 6: Interview Utilities

Teach:

- `heapq.heappush()`
- `heapq.heappop()`
- `heapq.heapify()`
- `math.inf`
- `math.gcd()`
- `divmod()`
- `itertools.accumulate()`
- `itertools.combinations()`

Applications:

- Top-K elements
- Priority queues
- Prefix sums
- Number manipulation
- Combination generation

### Module 7: Applied Interview Patterns

Connect Python language features to:

- Two pointers
- Sliding windows
- Prefix sums
- BFS
- DFS
- Binary search
- Heap-based selection
- Hash-map grouping

These are short reinforcement exercises, not a replacement for algorithm practice.

---

## 6. Exercise Data Model

Exercise definitions must be separate from the Go application.

Do not hardcode exercises into UI components.

Suggested exercise layout:

```text
exercises/
  iteration/
    enumerate-01/
      exercise.toml
      starter.py
      test_exercise.py
      solution.py
      explanation.md
```

### Example Metadata

```toml
id = "enumerate-01"
title = "Find an Index"
module = "iteration"
difficulty = "beginner"
type = "syntax"
order = 1

concepts = ["enumerate", "iteration"]

description = """
Find the index of the first occurrence
of a target value using enumerate().
Return -1 if the value is absent.
"""

objective = "Understand how enumerate() yields indices and values."

hints = [
  "Python can provide an index and value in the same loop.",
  "Try iterating with enumerate(nums).",
  "Return the index when the value matches the target."
]
```

### Starter Code

```python
def find_index(nums, target):
    # TODO: Use enumerate()
    pass
```

### Tests

```python
from starter import find_index

def test_target_present():
    assert find_index([10, 20, 30], 20) == 1

def test_target_missing():
    assert find_index([10, 20, 30], 50) == -1

def test_empty_list():
    assert find_index([], 5) == -1
```

The actual runner must execute the learner's working copy rather than the immutable starter template.

### Reference Solutions

Store reference solutions separately.

They should never overwrite user work unless an explicit reset is requested.

### Test Philosophy

Tests should focus primarily on correctness.

For syntax drills that specifically require a feature such as `enumerate()`, clearly communicate the requirement.

Where appropriate, include instructional checks that verify use of a designated API.

Do not implement brittle source-string checks that merely look for function names inside comments.

Keep feature-specific checks narrow, documented, and explainable.

---

## 7. Python Execution Engine

Implement a standalone Go package responsible for running exercise tests.

### Requirements

- Discover the configured Python interpreter.
- Check pytest availability.
- Execute each exercise's tests.
- Capture stdout and stderr.
- Capture exit codes.
- Enforce configurable timeouts.
- Support cancellation.
- Prevent indefinite UI blocking.
- Return structured test results.

Prefer pytest's machine-readable output where practical, or implement a small supported adapter that returns structured results without relying on fragile terminal-output parsing.

### Runtime Safety

Assume locally bundled, trusted curriculum exercises.

Executing arbitrary Python is not secure solely because it uses subprocesses.

Use controlled working directories and execution timeouts.

Do not claim to sandbox Python unless OS-level isolation is implemented.

Document the security implications of downloading third-party exercises.

Never build execution commands using shell-string concatenation.

### Background Execution

Python execution must not freeze the TUI.

Use Bubble Tea commands to perform subprocess work asynchronously and return typed messages to the update loop.

The interface should display running, completed, failed, and timed-out states.

Prevent concurrent test runs for the same exercise unless explicitly supported.

### Python Dependencies

Prefer a dedicated Python virtual environment for pytest.

Support configurable Python interpreter paths.

Avoid silently installing dependencies or modifying the user's global Python environment.

The `doctor` command should explain how to install missing requirements.

---

## 8. Editor Integration

The MVP should use an external code editor rather than implementing an embedded editor.

Resolve the editor in this order:

1. `$VISUAL`
2. `$EDITOR`
3. A documented fallback when available

Support Neovim, Vim, and other configurable editors.

### Editing Workflow

1. Select an exercise.
2. Press `e`.
3. Temporarily hand terminal control to the editor.
4. Open the learner's Python solution.
5. Let the learner save and exit.
6. Restore the TUI.
7. Allow the learner to run tests.

Use Bubble Tea v2's supported external-process mechanism.

Handle terminal restoration, editor failures, and interruption correctly.

Editor commands may contain arguments, so implement command parsing safely.

Do not execute arbitrary shell command strings.

### Working Files

Copy the starter implementation into a user-specific workspace the first time an exercise is opened.

Preserve user changes across application restarts.

Never edit the packaged starter file directly.

---

## 9. Progress Persistence

Use a JSON file in the platform-appropriate user data directory.

Track:

- Completed exercises
- Attempt counts
- Last attempt timestamp
- Last successful completion timestamp
- Overall completion percentage
- Completion percentage per module

An exercise is completed when its required tests pass.

Do not overwrite completion history merely because the user reopens an exercise.

Use atomic file writes.

Keep progress independent from source exercise files.

Progress must persist across application restarts.

---

## 10. CLI Commands

Implement:

```bash
leetkit
leetkit list
leetkit run enumerate-01
leetkit check enumerate-01
leetkit progress
leetkit reset enumerate-01
leetkit doctor
```

### Expected Behavior

`leetkit` opens the dashboard.

`leetkit list` prints available modules and exercises.

`leetkit run <id>` opens the selected exercise in the interactive interface.

`leetkit check <id>` executes tests without opening the TUI.

`leetkit progress` displays completion statistics.

`leetkit reset <id>` restores the original starter code and clears relevant progress after user confirmation. Provide an explicit noninteractive confirmation flag for scripts.

`leetkit doctor` checks runtime requirements and configuration.

All commands must report actionable errors.

---

## 11. Visual Design

Use Tokyo Night as the default theme.

Suggested colors:

```text
Background:    #1a1b26
Secondary BG:  #24283b
Foreground:    #c0caf5
Blue:          #7aa2f7
Purple:        #bb9af7
Green:         #9ece6a
Yellow:        #e0af68
Red:           #f7768e
Muted:         #565f89
```

### Design Principles

- Clean, minimal layout
- Consistent typography and borders
- Strong keyboard focus indicators
- Clear completed/incomplete states
- Readable code previews
- Subtle success/failure feedback
- No unnecessary animations
- Responsive terminal resizing

Use Bubbles components for lists, viewports, progress indicators, help, and text inputs where appropriate.

Avoid recreating functionality already supported by Charm.

Handle terminals without true-color support gracefully.

---

## 12. Go Project Structure

Use idiomatic Go package organization.

Suggested structure:

```text
leetkit/
├── go.mod
├── go.sum
├── main.go
├── README.md
├── LICENSE
├── Makefile
├── .gitignore
│
├── cmd/
│   └── leetkit/
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   │
│   ├── ui/
│   │   ├── dashboard.go
│   │   ├── exercises.go
│   │   ├── details.go
│   │   ├── results.go
│   │   ├── progress.go
│   │   └── theme.go
│   │
│   ├── exercise/
│   │   ├── model.go
│   │   ├── loader.go
│   │   └── validation.go
│   │
│   ├── runner/
│   │   ├── runner.go
│   │   ├── python.go
│   │   └── results.go
│   │
│   ├── editor/
│   │   └── editor.go
│   │
│   ├── progress/
│   │   └── store.go
│   │
│   └── config/
│       └── config.go
│
├── exercises/
│   ├── iteration/
│   ├── hashing/
│   ├── collections/
│   ├── sorting/
│   ├── pythonic/
│   ├── utilities/
│   └── patterns/
│
└── tests/
```

Treat this as guidance rather than a mandatory one-file-per-concept layout.

Prefer a small number of cohesive packages over excessive abstraction.

Keep business logic separate from UI rendering.

Ensure the exercise loader, test runner, and progress store can be tested without launching an interactive terminal.

### Exercise Distribution

Bundled exercises must work with installed release binaries, not just from the repository root.

Consider Go's `embed` package for distributing official exercise templates.

Extract or copy working files into a user data directory as needed.

---

## 13. MVP Scope

Do not attempt to implement the full 45–50 exercise curriculum immediately.

The MVP must include exactly 15 well-designed exercises.

### MVP Exercise Set

| ID | Topic |
|---|---|
| 01 | `enumerate()` |
| 02 | `zip()` |
| 03 | `dict.get()` |
| 04 | Sets |
| 05 | `Counter` |
| 06 | `defaultdict` |
| 07 | `deque` |
| 08 | List comprehensions |
| 09 | Dictionary comprehensions |
| 10 | `sorted(..., key=...)` |
| 11 | `lambda` |
| 12 | Tuple unpacking |
| 13 | `heapq` |
| 14 | `bisect` |
| 15 | `itertools.accumulate()` |

Each exercise must include:

- Metadata
- Starter code
- Tests
- Hints
- Reference solution
- Explanation
- Interview relevance

Favor high-quality exercise content over increasing exercise count.

---

## 14. Development Milestones

### Phase 1 — Foundation

Implement:

- Go project initialization
- Bubble Tea v2 application
- Main dashboard
- Keyboard navigation
- Exercise metadata loader
- Basic styling
- CLI entry point

Acceptance criteria:

- `go run .` launches the TUI.
- Users can browse bundled exercises.
- Window resizing works.
- Keyboard navigation works.
- Application exits cleanly.

### Phase 2 — Exercise Engine

Implement:

- Python interpreter discovery
- pytest execution
- Working directory management
- External editor integration
- Structured test results
- Execution timeouts
- Error handling

Acceptance criteria:

- A user can open an exercise.
- Edit its Python code.
- Return to LeetKit.
- Run the exercise's tests.
- See correct pass/fail feedback.

### Phase 3 — Learning Features

Implement:

- Exercise hints
- Explanations
- Completion tracking
- Progress persistence
- Progress dashboard
- Exercise search
- Keyboard help

Acceptance criteria:

- Progress survives restarts.
- Completed exercises remain completed.
- Hints work correctly.
- Reference explanations are accessible.
- Search works across modules.

### Phase 4 — Curriculum

Implement all 15 MVP exercises.

Each must have meaningful tests and educational explanations.

Verify the bundled reference solutions.

### Phase 5 — Polish and Release Readiness

Implement:

- Responsive layouts
- Better error presentation
- Improved startup behavior
- README documentation
- Installation instructions
- Contribution documentation
- Release build configuration

Do not add cloud synchronization, authentication, AI tutoring, achievements, or embedded code editing in the MVP.

These can be evaluated later.

---

## 15. Testing and Quality

The project must pass:

```bash
go fmt ./...
go vet ./...
go test ./...
```

Where practical, also use:

```bash
go test -race ./...
```

Testing should cover:

- TOML parsing
- Invalid exercise definitions
- Exercise discovery
- Duplicate exercise IDs
- Subprocess failures
- Execution timeouts
- Python syntax errors
- Editor launch failures
- Progress persistence
- Progress recovery
- Terminal resizing
- TUI event transitions
- CLI command handling

Also execute every exercise's reference solution against its test suite.

Do not consider the MVP complete if exercise content or the execution workflow is broken.

### Important Edge Cases

Handle:

- Missing Python installation
- Missing pytest
- Missing editor
- Corrupted progress file
- Interrupted Python execution
- Empty exercise directories
- Terminal too small
- Exercises containing syntax errors
- Invalid exercise IDs
- Unexpected subprocess termination

Display helpful recovery instructions.

---

## 16. Documentation

Create a clear README containing:

1. Project overview
2. Motivation
3. Screenshots or terminal demonstrations when available
4. Requirements
5. Installation
6. Quick start
7. Keyboard shortcuts
8. Available commands
9. Exercise structure
10. How to contribute exercises
11. Development instructions
12. Project roadmap

The README should emphasize:

**LeetKit helps developers learn Python's most useful language features for coding interviews.**

Clarify that the project is independent and not affiliated with LeetCode.

Do not claim official LeetCode integration.

---

## 17. Out of Scope

Do not implement the following in the MVP:

- Full LeetCode problem sets
- LeetCode API integration
- Cloud execution
- User accounts
- Leaderboards
- Streak systems
- AI-generated hints
- Online synchronization
- Embedded code editor
- Multiple programming languages
- Full algorithm courses
- Remote exercise downloads

The project is intentionally Python-only.

Its main differentiator is mastering Python's interview-relevant language features.

---

## 18. Instructions for Claude Code

Before implementation:

1. Inspect the current repository.
2. Read any existing project instructions or `CLAUDE.md`.
3. Review the current official Bubble Tea v2, Bubbles v2, and Lip Gloss v2 documentation.
4. Design a maintainable minimal architecture.
5. Identify dependencies and implementation milestones.

Then implement the application incrementally.

### Engineering Expectations

- Write idiomatic Go.
- Use current compatible library APIs.
- Prefer standard-library solutions when practical.
- Keep UI state management predictable.
- Avoid unnecessary dependencies.
- Keep exercise content data-driven.
- Use meaningful error handling.
- Do not silently swallow errors.
- Avoid speculative abstractions.
- Maintain separation between UI, execution, and persistence.
- Verify functionality before moving to the next milestone.

### Execution Strategy

Build the smallest complete vertical slice first:

**Browse one exercise → open editor → write Python → run tests → see feedback → save progress.**

Once that works, implement the remaining dashboard features and expand the curriculum.

Do not build a visually impressive but nonfunctional prototype.

The primary acceptance criterion is an end-to-end learning workflow.

### Final Deliverables

Produce:

- Working Go CLI/TUI
- Functional Python execution engine
- 15 complete educational exercises
- Persistent local progress
- External editor integration
- Automated tests
- README and contribution documentation
- Clean, maintainable source code

Run the relevant checks and report any remaining limitations.

## Final Product Vision

LeetKit should become a lightweight, enjoyable tool developers can use for 10–15 minutes a day to sharpen their Python skills before coding interviews.

It should feel like **Rustlings meets a polished Charm TUI**, with exercises specifically designed around Python's most useful interview features.

The ultimate goal isn't to teach users how to solve every algorithm problem.

It's to make Python syntax and standard-library tools feel second nature, so developers can focus on the algorithms themselves.
