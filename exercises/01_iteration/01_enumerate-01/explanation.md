# enumerate()

`enumerate(iterable, start=0)` yields `(index, value)` pairs.

```python
for i, num in enumerate(nums):
    ...
```

It replaces the index-based loop:

```python
for i in range(len(nums)):
    num = nums[i]
```

The `enumerate` version is shorter, works with any iterable (not only lists),
and removes a class of off-by-one bugs.

## Useful variations

- `enumerate(nums, start=1)` counts from 1, handy for "position" style output.
- `{num: i for i, num in enumerate(nums)}` builds a value-to-index map.
- `next((i for i, x in enumerate(nums) if x == target), -1)` finds the first
  match in one expression.

## Interview relevance

Reach for `enumerate` whenever you need the position of an element while
iterating: Two Sum (store `value -> index`), finding the first or last index
that satisfies a condition, or comparing an element with its neighbour by index.
Interviewers notice idiomatic loops; `range(len(nums))` reads as unfamiliarity
with Python.
