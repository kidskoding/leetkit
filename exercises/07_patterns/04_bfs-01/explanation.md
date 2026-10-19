# BFS with collections.deque

Breadth-first search explores a region layer by layer from a FIFO queue.
`collections.deque` is the right queue: `append` and `popleft` are both O(1),
while `list.pop(0)` shifts every element.

```python
from collections import deque

seen = {(r, c)}
queue = deque([(r, c)])
while queue:
    cr, cc = queue.popleft()
    for nr, nc in ((cr + 1, cc), (cr - 1, cc), (cr, cc + 1), (cr, cc - 1)):
        if 0 <= nr < rows and 0 <= nc < cols and grid[nr][nc] == "1" and (nr, nc) not in seen:
            seen.add((nr, nc))
            queue.append((nr, nc))
```

Tuples are hashable, so a `set` of `(row, col)` pairs is a natural visited
marker that leaves the input untouched. Chained comparisons like
`0 <= nr < rows` keep the bounds check readable. Mark a cell when you enqueue
it, not when you pop it, or the same cell can enter the queue several times.

## Useful variations

- `for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):` keeps directions in one constant.
- `queue.append((nr, nc, dist + 1))` carries distance for shortest-path questions.
- `for _ in range(len(queue)):` processes one full level per outer step (rotting oranges).
- Overwriting `grid[r][c] = "0"` saves the set, but only when mutating input is allowed.

## Interview relevance

Number of islands, max area of island, rotting oranges, walls and gates, and
shortest path in a binary matrix are all grid BFS. Interviewers watch for an
O(1) queue, correct bounds checks, and whether you ask before mutating the input.
