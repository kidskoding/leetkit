# Nested comprehensions

A comprehension can appear inside another comprehension's output expression.
That builds a list of lists, one fresh inner list per outer iteration:

```python
grid = [[0 for _ in range(cols)] for _ in range(rows)]
```

This avoids the classic aliasing bug:

```python
grid = [[0] * cols] * rows   # every row is the SAME list
grid[0][0] = 1               # changes column 0 of every row
```

A single comprehension can also have several `for` clauses. They read in the
same order as the equivalent nested loops, outer loop first:

```python
flat = [x for row in matrix for x in row]

# same as
flat = []
for row in matrix:
    for x in row:
        flat.append(x)
```

The two shapes do different things: `[[...] for ...]` keeps the structure,
`[... for ... for ...]` flattens it.

## Useful variations

- `[[0] * cols for _ in range(rows)]` is also safe; only the outer `*` aliases.
- `[[matrix[r][c] for r in range(rows)] for c in range(cols)]` transposes.
- `[(r, c) for r in range(rows) for c in range(cols)]` lists every cell coordinate.
- `list(itertools.chain.from_iterable(matrix))` is another way to flatten.

## Interview relevance

DP tables, visited arrays and game boards all need a correctly allocated grid,
and interviewers do notice `[[0] * n] * m`. Flattening shows up when you need
to sort, count or take the median of every value in a matrix. Beyond two
levels, nested comprehensions get hard to read; switch to explicit loops.
