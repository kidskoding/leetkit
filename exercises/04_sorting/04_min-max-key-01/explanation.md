# min() and max() with key=

`min` and `max` take the same `key=` argument as `sorted`. They return the
original element whose key is smallest or largest, in a single O(n) pass.

```python
points = [(3, 4), (1, -1), (5, 0)]
min(points, key=lambda p: p[0] * p[0] + p[1] * p[1])
# (1, -1)
```

Compare with the manual version:

```python
best = None
for p in points:
    if best is None or p[0] ** 2 + p[1] ** 2 < best[0] ** 2 + best[1] ** 2:
        best = p
```

On ties, both `min` and `max` return the first element with the winning key.
On an empty iterable they raise `ValueError`, unless you pass `default=`.

## Useful variations

- `max(words, key=len)` returns the longest word (first one on ties).
- `max(counts, key=counts.get)` returns the dict key with the largest value.
- `min(nums, default=0)` returns 0 instead of raising on an empty list.
- `sorted(points, key=dist)[:k]` or `heapq.nsmallest(k, points, key=dist)`
  when you need the best k, not just one.

## Interview relevance

Choosing the best candidate by some score shows up in K Closest Points, Most
Frequent Element, Longest Word and greedy algorithms. `min(..., key=...)` is
linear, avoids an unnecessary sort, and makes the scoring rule explicit.
Knowing that ties go to the first element and that `default=` avoids a crash
on empty input covers the edge cases interviewers ask about.
