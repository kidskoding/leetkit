# heapq.heappush and heapq.heappop

`heapq` turns a plain list into a binary min-heap. The smallest element is
always at `heap[0]`.

```python
import heapq

heap = []
heapq.heappush(heap, 5)
heapq.heappush(heap, 1)
heapq.heappop(heap)  # 1, the smallest
```

For "kth largest", keep a min-heap of size k. Whenever it grows past k, pop the
smallest. What remains are the k largest values, and `heap[0]` is the kth largest.

```python
for num in nums:
    heapq.heappush(heap, num)
    if len(heap) > k:
        heapq.heappop(heap)
```

## Useful variations

- `heapq.heappushpop(heap, num)` pushes then pops in one faster call.
- `heapq.nlargest(k, nums)` returns the k largest as a sorted list.
- Push `-num` to simulate a max-heap, since `heapq` only provides a min-heap.

## Interview relevance

Top-k problems (Kth Largest Element, K Closest Points, Top K Frequent Elements)
expect the O(n log k) heap solution. Knowing that `heapq` is a min-heap, and
that a min-heap of size k keeps the largest values, is the key insight.
