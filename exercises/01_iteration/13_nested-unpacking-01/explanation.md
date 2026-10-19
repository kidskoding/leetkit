# Nested tuple unpacking

A for-loop target can mirror the shape of each item, so nested tuples unpack in
one step.

```python
points = [((1, 1), "cafe"), ((5, 2), "park")]
for (x, y), label in points:
    print(label, x, y)
```

The same works in plain assignments and comprehensions:

```python
(x_min, y_min), (x_max, y_max) = box
[label for (x, y), label in points if x <= x_max]
```

Compare with indexing, where every reader has to decode the positions:

```python
for p in points:
    if p[0][0] <= box[1][0] and p[0][1] <= box[1][1]:
        result.append(p[1])
```

If the shape does not match (a 3-tuple where 2 were expected), Python raises
ValueError immediately, which catches malformed data early.

## Useful variations

- `for i, (key, value) in enumerate(d.items()):` combines an index with a pair.
- `dist, (r, c) = heapq.heappop(heap)` unpacks a priority-queue entry.
- `for (r, c), steps in queue:` reads naturally in BFS over a grid.
- `(first, *rest), last = ...` mixes nesting with star unpacking.

## Interview relevance

BFS queues hold `((r, c), dist)`, Dijkstra heaps hold `(cost, (node, prev))`, and
interval or point problems hold pairs of pairs. Unpacking the whole shape in the
loop header gives every piece a name, which keeps the logic readable and avoids
mixing up `[0]` and `[1]` while talking through the solution.
