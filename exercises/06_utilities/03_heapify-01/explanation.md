# heapq.heapify()

`heapq.heapify(lst)` rearranges a list into a min-heap in place, in O(n) time.
It returns `None`.

```python
import heapq

heap = list(prices)
heapq.heapify(heap)
smallest = [heapq.heappop(heap) for _ in range(min(k, len(heap)))]
```

Because it works in place, heapifying the argument directly would reorder the
caller's list:

```python
heapq.heapify(prices)  # prices is now in heap order, not the order you were given
```

Copy first with `list(prices)` or `prices[:]` so the input stays untouched.

## Useful variations

- `heapq.nsmallest(k, prices)` does the same job in one call and never mutates the input.
- `heap = heapq.heapify(prices)` is a bug: `heap` is `None`.
- Heapify a list of tuples, such as `[(dist, point) for ...]`, to build a priority queue in one step.

## Interview relevance

When all items are known up front (K Closest Points, Last Stone Weight),
`heapify` builds the heap faster than repeated pushes. Copying before mutating
shows care with shared data, which interviewers look for.
