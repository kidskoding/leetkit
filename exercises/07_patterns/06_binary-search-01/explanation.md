# Binary search on a predicate

When a yes/no check is False for a while and then True forever, binary search
finds the flip point in O(log n) checks. Keep a half-open range `[lo, hi)` that
always contains the answer, and shrink it until it is empty.

```python
lo, hi = 1, n + 1
while lo < hi:
    mid = (lo + hi) // 2
    if is_bad(mid):
        hi = mid
    else:
        lo = mid + 1
return lo if lo <= n else -1
```

Python integers never overflow, so `(lo + hi) // 2` is safe without the
`lo + (hi - lo) // 2` trick. Because `hi = mid` (not `mid - 1`) and
`lo = mid + 1`, the range always shrinks and the loop always ends.

## Useful variations

- `bisect.bisect_left(range(1, n + 1), True, key=is_bad)` does the same search (Python 3.10+).
- `bisect_left(sorted_list, x)` finds the first index with value `>= x` in a sorted list.
- Integer square root: search for the first `m` with `m * m > x`, then answer `m - 1`.
- `math.isqrt(x)` exists, but interviewers usually want the search written out.

## Interview relevance

First bad version, search insert position, Koko eating bananas, capacity to
ship packages, and split array largest sum are all "find the first True".
Interviewers care less about the code than about a loop invariant you can state
and boundaries that do not loop forever.
