# heapq.merge()

`heapq.merge(*iterables)` takes several already-sorted iterables and yields
their items as one sorted stream. Internally it keeps a heap of the current
front item of each input, so it runs in O(N log k) for k inputs.

```python
import heapq

list(heapq.merge([1, 4, 9], [2, 3, 10], [5]))  # [1, 2, 3, 4, 5, 9, 10]
```

The shortcut most people reach for first ignores the fact that the inputs are
already sorted:

```python
sorted(x for stream in streams for x in stream)  # O(N log N), builds everything
```

`merge` is lazy: it returns an iterator and only pulls the next item from each
input when needed, so it works on huge files or generators. Wrap it in `list()`
when you need a list.

## Useful variations

- `heapq.merge(a, b, key=len)` merges by a key; each input must be sorted by that same key.
- `heapq.merge(a, b, reverse=True)` merges inputs that are sorted in descending order.
- `next(heapq.merge(*streams))` gets the overall smallest item without merging the rest.
- For linked lists (LeetCode 23), push `(node.val, i, node)` tuples onto a heap yourself; `merge` needs iterables.

## Interview relevance

Merge k Sorted Lists, Kth Smallest in a Sorted Matrix, and external sorting all
use the k-way merge pattern. Knowing `heapq.merge` lets you solve the
iterable version in one line and explain the heap mechanics behind it.
