# heapq.nlargest() and heapq.nsmallest()

`heapq.nlargest(k, iterable, key=None)` returns the `k` largest items, largest
first. `heapq.nsmallest` is the mirror image. With `key=` you rank whole
records by one field and get the records back.

```python
import heapq

orders = [{"id": "a1", "total": 25}, {"id": "b2", "total": 90}, {"id": "c3", "total": 40}]
heapq.nlargest(2, orders, key=lambda o: o["total"])
# [{'id': 'b2', 'total': 90}, {'id': 'c3', 'total': 40}]
```

It behaves like `sorted(orders, key=..., reverse=True)[:k]`, including keeping
equal items in input order, but only keeps a heap of size `k`, so it costs
O(n log k) instead of O(n log n). If `k` is bigger than the input you simply
get everything, and `k = 0` gives `[]`.

## Useful variations

- `heapq.nsmallest(3, points, key=lambda p: p[0] ** 2 + p[1] ** 2)` finds the 3 closest points.
- `heapq.nlargest(k, counts, key=counts.get)` returns the k most frequent keys of a dict.
- `operator.itemgetter("total")` is a readable alternative to the lambda.
- For `k == 1`, `max(orders, key=...)` is simpler and faster.

## Interview relevance

Top K Frequent Elements, K Closest Points to Origin, and leaderboard-style
questions all reduce to "rank by a key, keep k". Using `nlargest` or `nsmallest`
shows you know the heap-based O(n log k) approach without hand-rolling it.
