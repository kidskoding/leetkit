# bisect.bisect_left

`bisect_left(a, x)` binary-searches the sorted list `a` and returns the first
index `i` such that `a[i] >= x`. That is exactly where `x` would be inserted
to keep `a` sorted, and if `x` is present it is the index of its leftmost copy.

```python
from bisect import bisect_left

bisect_left([1, 3, 5, 6], 5)   # 2  (found)
bisect_left([1, 3, 5, 6], 2)   # 1  (would go before 3)
bisect_left([1, 3, 5, 6], 7)   # 4  (past the end)
```

The hand-written version needs `lo`, `hi`, a midpoint and careful choice of
`<` versus `<=`. `bisect_left` runs in O(log n) and gets the boundaries right.

## Useful variations

- `i < len(a) and a[i] == x` after `i = bisect_left(a, x)` tests membership.
- `bisect_right(a, x)` returns the index just past the last copy of `x`.
- `insort(a, x)` inserts `x` in sorted position (O(n) because of the shift).
- `bisect_left(a, x, key=...)` searches by a key on Python 3.10+.

## Interview relevance

Any "first index where value >= x" question is a `bisect_left` call: Search
Insert Position, the patience-sorting step of Longest Increasing Subsequence,
and timestamp lookups. Using it frees you to focus on the problem; if the
interviewer wants the manual loop, knowing what `bisect_left` returns tells you
exactly which invariant to implement.
