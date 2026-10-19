# Top-k selection with heapq.nsmallest

`heapq.nsmallest(k, iterable, key=...)` returns the `k` smallest items, smallest
first, using a heap of size `k` internally. It is stable: items with equal keys
come out in input order.

```python
import heapq

closest = heapq.nsmallest(k, points, key=lambda p: p[0] * p[0] + p[1] * p[1])
```

The hand-rolled version pushes negated distances into a max-heap and pops
whenever it grows past `k`. That is worth knowing, but in Python the one-liner
expresses the same O(n log k) idea. Comparing squared distances is enough to
rank points, avoids `math.sqrt`, and keeps every key an exact integer.

## Useful variations

- `heapq.nlargest(k, counts, key=counts.get)` gives the top k frequent keys.
- `heapq.heappush(heap, (-dist, i, p))` plus `heappop` when `len(heap) > k` is the manual max-heap.
- `sorted(points, key=...)[:k]` is fine when `k` is close to `n`.
- `math.dist(p, (0, 0))` gives the real distance when you need the value itself.

## Interview relevance

K closest points, top k frequent elements, kth largest element, and merging k
sorted lists all lean on `heapq`. Knowing `nsmallest`/`nlargest` with a `key`
lets you finish quickly, then discuss the manual heap and complexity if asked.
