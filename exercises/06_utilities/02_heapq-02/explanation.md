# Priority queues with heapq and tuples

Python compares tuples element by element, so a heap of tuples is ordered by
the first field, then the second, and so on.

```python
import heapq

heap = []
for i, (name, priority) in enumerate(tasks):
    heapq.heappush(heap, (priority, i, name))

while heap:
    priority, _, name = heapq.heappop(heap)
```

The arrival index `i` does two jobs. It breaks ties so equal priorities come
out in arrival order, and because every index is unique, Python never needs to
compare `name`. Without it, a tie would compare the payloads, which is wrong
for strings and raises `TypeError` for objects such as dicts or nodes.

## Useful variations

- `heapq.heapify(list_of_tuples)` builds the queue in O(n) when all items are known up front.
- Push `(-priority, i, item)` when larger numbers mean more urgent.
- `itertools.count()` gives a running tie-breaker: `heappush(heap, (p, next(counter), item))`.

## Interview relevance

Dijkstra pushes `(distance, node)`, Merge K Sorted Lists pushes
`(value, list_index, node)`, and Task Scheduler style problems push
`(time, id, task)`. Remembering the tie-breaker field avoids a classic
`TypeError: '<' not supported` crash in the middle of an interview.
