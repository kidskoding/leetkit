# Counting a range with bisect

In a sorted list, all values between `lo` and `hi` sit in one contiguous
block. Two binary searches find its edges:

```python
from bisect import bisect_left, bisect_right

start = bisect_left(prices, lo)    # first index with value >= lo
end = bisect_right(prices, hi)     # first index with value > hi
count = end - start
```

`bisect_left` lands before any copies of `lo`, so they are included.
`bisect_right` lands after any copies of `hi`, so they are included too. That
pairing gives an inclusive range. When `lo > hi` the difference can be
negative, so clamp it with `max(0, ...)`.

The linear version, `sum(lo <= p <= hi for p in prices)`, is O(n) per query;
the bisect version is O(log n).

## Useful variations

- `bisect_left(a, hi) - bisect_left(a, lo)` counts the half-open range
  `lo <= p < hi`.
- `bisect_right(a, x) - bisect_left(a, x)` counts the copies of `x`.
- `prices[bisect_left(prices, lo):bisect_right(prices, hi)]` returns the
  values themselves instead of the count.

## Interview relevance

"Sort once, then answer many range queries" is a common follow-up to brute
force solutions: counting events in a time window, scores in a band, or pairs
whose sum falls in a range. Knowing which bisect variant gives an inclusive or
exclusive boundary is what keeps these answers free of off-by-one errors.
