# Grid neighbours with a DIRECTIONS list

Encode "up, down, left, right" as row and column offsets once, then loop over
them. One bounds check covers every direction.

```python
DIRECTIONS = [(-1, 0), (1, 0), (0, -1), (0, 1)]

for dr, dc in DIRECTIONS:
    nr, nc = r + dr, c + dc
    if 0 <= nr < rows and 0 <= nc < cols:
        ...  # heights[nr][nc] is a real neighbour
```

The alternative is four near-identical blocks, each with its own bounds test:

```python
if r > 0 and heights[r - 1][c] >= h: ...
if r < rows - 1 and heights[r + 1][c] >= h: ...
if c > 0 and heights[r][c - 1] >= h: ...
if c < cols - 1 and heights[r][c + 1] >= h: ...
```

A copy-paste slip in one of those (`r - 1` where `c - 1` was meant) is a common
grid bug and hard to spot.

## Useful variations

- Add `(-1, -1), (-1, 1), (1, -1), (1, 1)` for 8-directional problems such as
  Game of Life.
- `all(h > grid[r + dr][c + dc] for dr, dc in DIRECTIONS if in_bounds(...))`
  turns a whole neighbour check into one expression; `all()` of nothing is True.
- In BFS, `for dr, dc in DIRECTIONS:` builds the next frontier of `(nr, nc)` cells.
- A small `def in_bounds(r, c): return 0 <= r < rows and 0 <= c < cols` keeps
  the condition readable.

## Interview relevance

Number of Islands, Rotting Oranges, Walls and Gates, Word Search, Game of Life,
and Pacific Atlantic Water Flow all start with "for each neighbour of this
cell". Writing the DIRECTIONS list first is the standard opening move: it keeps
the traversal code short and lets you switch between 4 and 8 neighbours by
editing one line.
