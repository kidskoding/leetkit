# Longest Increasing Subsequence with bisect_left

Keep a list `tails` where `tails[k]` is the smallest last value of any
increasing subsequence of length `k + 1` seen so far. That list is always
sorted, so each new value finds its place with one binary search.

```python
from bisect import bisect_left

tails = []
for x in nums:
    i = bisect_left(tails, x)
    if i == len(tails):
        tails.append(x)
    else:
        tails[i] = x
return len(tails)
```

If `x` is larger than every tail it extends the longest subsequence. Otherwise
it replaces the first tail that is `>= x`, giving that length a smaller,
easier-to-extend ending. The O(n^2) DP compares every pair; this version is
O(n log n). Note that `tails` is not itself a valid subsequence, only its length is meaningful.

## Useful variations

- `bisect_left` gives strictly increasing; `bisect_right` gives non-decreasing (equal values allowed).
- Store the index where each value landed to rebuild an actual subsequence afterwards.
- Russian Doll Envelopes: sort by `(width, -height)`, then run this on the heights.
- Longest decreasing subsequence: run it on `[-x for x in nums]`.

## Interview relevance

LIS is a staple, and interviewers often accept the O(n^2) DP first and then ask
for better. Recognising the tails invariant and writing it with `bisect_left`
turns that follow-up into a few lines, and the same trick solves the
envelope-nesting and box-stacking variants.
