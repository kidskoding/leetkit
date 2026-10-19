# bisect.insort

`insort(a, x)` inserts `x` into the sorted list `a` at the position that keeps
it sorted. It finds the spot with a binary search, then shifts the later
elements over, so the list never needs re-sorting.

```python
from bisect import insort

window = []
for value in [5, 15, 1, 3]:
    insort(window, value)
# window goes [5] -> [5, 15] -> [1, 5, 15] -> [1, 3, 5, 15]
```

The clunky version appends and calls `window.sort()` every time, which is
O(n log n) per value. `insort` is O(log n) to search plus O(n) to shift, and it
states the intent: this list stays sorted.

## Useful variations

- `insort_left` and `insort_right` decide whether equal values go before or after existing copies; `insort` is `insort_right`.
- `insort(a, x, key=...)` (Python 3.10+) keeps records sorted by a field; `a` must already be in key order.
- `a.pop(bisect_left(a, x))` removes one copy of `x`, which is how a sliding-window median drops old values.
- Two heaps give O(log n) per value when the stream is huge; `insort` is simpler for moderate sizes.

## Interview relevance

Find Median from Data Stream and Sliding Window Median both need "the middle
of everything seen so far". A sorted list kept with `insort` gets a correct
solution on the board fast, and knowing its O(n) shift cost lets you explain
when to upgrade to the two-heap design.
