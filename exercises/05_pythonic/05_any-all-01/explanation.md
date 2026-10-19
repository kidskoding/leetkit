# any() and all()

`any(iterable)` is True if at least one item is truthy. `all(iterable)` is True
if every item is truthy. Pair them with a generator expression to turn a rule
into one line:

```python
has_digit = any(ch.isdigit() for ch in password)
no_spaces = all(not ch.isspace() for ch in password)
```

Both short-circuit: `any` stops at the first True, `all` stops at the first
False. The loop they replace looks like this:

```python
has_digit = False
for ch in password:
    if ch.isdigit():
        has_digit = True
        break
```

Edge cases are worth remembering: `any([])` is False and `all([])` is True.
An empty password passes the whitespace rule on its own, which is why the
length rule has to reject it.

## Useful variations

- `all(len(row) == len(grid[0]) for row in grid)` checks a rectangular grid.
- `any(x < 0 for x in nums)` detects a negative value.
- `all(a <= b for a, b in zip(nums, nums[1:]))` checks that a list is sorted.
- `not any(...)` reads as "none of"; it is equivalent to `all(not ...)`.

## Interview relevance

Validity checks appear in Valid Sudoku, "is the matrix Toeplitz", "is this
array monotonic", and many simulation problems. Writing each condition as an
`any` or `all` keeps the logic declarative and easy to test out loud, and the
short-circuiting means you do not pay for checks you do not need.
