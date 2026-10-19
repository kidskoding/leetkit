# range(start, stop, step)

`range(start, stop, step)` produces integers from `start` up to, but not
including, `stop`, moving by `step` each time.

```python
range(0, 10, 3)     # 0, 3, 6, 9
range(9, -1, -3)    # 9, 6, 3, 0
```

With a negative step, `range` counts down. Because `stop` is exclusive, the
way to include index `0` is to stop at `-1`:

```python
for i in range(len(nums) - 1, -1, -1):
    ...  # last index down to 0
```

A common bug is `range(len(nums) - 1, 0, -1)`, which silently skips index 0.

## Useful variations

- `range(n)` is `range(0, n, 1)`; `range(1, n)` skips the first index.
- `range(0, len(s), 2)` walks even positions; `range(1, len(s), 2)` the odd ones.
- `reversed(range(n))` is an alternative to `range(n - 1, -1, -1)` that is
  harder to get wrong.

## Interview relevance

Right-to-left scans power monotonic-stack problems (Daily Temperatures, Next
Greater Element), suffix products, and in-place array rewrites. Stride loops
appear in chunking and sampling tasks. Being fluent with `range`'s three
arguments, and its exclusive stop, avoids off-by-one errors under pressure.
