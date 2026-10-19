# deque for BFS

`collections.deque` is a double-ended queue with O(1) `append` and `popleft`.
That makes it the right container for breadth-first search.

```python
from collections import deque

queue = deque([(src, 0)])
visited = {src}
while queue:
    node, dist = queue.popleft()
    if node == dst:
        return dist
    for nxt in graph[node]:
        if nxt not in visited:
            visited.add(nxt)
            queue.append((nxt, dist + 1))
return -1
```

A plain list works but `queue.pop(0)` moves every remaining element, so each
pop is O(n) and the whole search can become O(n^2).

## Useful variations

- Process one level at a time with `for _ in range(len(queue)):` and a separate `steps` counter.
- Mark nodes visited when you enqueue them, not when you pop them, to avoid duplicates in the queue.
- For grids, enqueue `(row, col)` tuples and loop over `((1, 0), (-1, 0), (0, 1), (0, -1))`.

## Interview relevance

Any "minimum number of steps" question on an unweighted graph or grid is BFS:
Word Ladder, Rotting Oranges, Shortest Path in Binary Matrix, Open the Lock.
Using `deque.popleft()` is expected; reaching for `list.pop(0)` is a common
red flag.
