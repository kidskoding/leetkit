# List comprehensions

A list comprehension builds a new list from an iterable, optionally
transforming and filtering each item.

```python
squares = [n * n for n in nums if n % 2 == 0]
```

Read it left to right as "n squared, for each n in nums, if n is even". It
replaces the four-line accumulator loop:

```python
squares = []
for n in nums:
    if n % 2 == 0:
        squares.append(n * n)
```

The comprehension is shorter, has no mutable accumulator to get wrong, and is
usually a little faster.

## Useful variations

- `[x for x in nums if x > 0]` filters without transforming.
- `[s.strip().lower() for s in words]` transforms without filtering.
- `[i for i, x in enumerate(nums) if x == target]` collects matching indices.
- `[x if x > 0 else 0 for x in nums]` puts a conditional expression in the
  output position (note the `if` moves to the front when there is an `else`).

## Interview relevance

Comprehensions show up constantly as glue: building neighbour lists in a grid
BFS, filtering candidates, mapping characters to codes. Interviewers read them
as fluency. Keep them to one clear idea; if a comprehension needs nested
conditions and several clauses, a plain loop is easier to explain out loud.
