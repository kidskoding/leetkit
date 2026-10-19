# Transposing with zip(*matrix)

`*matrix` passes each row to `zip` as a separate argument, so `zip` walks all
rows in lockstep and yields one tuple per column.

```python
matrix = [[1, 2, 3],
          [4, 5, 6]]
list(zip(*matrix))   # [(1, 4), (2, 5), (3, 6)]
```

That is the transpose. Rotations follow from combining it with a reversal:

```python
[list(r) for r in zip(*reversed(matrix))]   # 90 degrees clockwise
[list(r) for r in zip(*matrix)][::-1]       # 90 degrees counter-clockwise
```

Compare with the index version:

```python
rows, cols = len(matrix), len(matrix[0])
out = [[matrix[rows - 1 - r][c] for r in range(rows)] for c in range(cols)]
```

## Useful variations

- `[sum(col) for col in zip(*matrix)]` gives column totals.
- `any(len(set(col)) < len(col) for col in zip(*board))` checks columns for
  duplicates (Valid Sudoku).
- `zip` yields tuples, so wrap with `list(...)` when the caller expects lists.
- `zip(*matrix)` on `[]` or `[[]]` yields nothing, which handles empty grids.

## Interview relevance

Rotate Image, Transpose Matrix, Spiral Matrix (peel a row, rotate the rest),
and column checks in Sudoku or Tic-Tac-Toe all use this idiom. It replaces
nested index loops, where swapping `i` and `j` is easy to get wrong.
