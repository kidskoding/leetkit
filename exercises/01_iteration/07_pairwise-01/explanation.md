# Adjacent pairs with zip(xs, xs[1:])

Zipping a list with itself shifted by one gives every neighbouring pair.

```python
temps = [12, 15, 9, 10]
list(zip(temps, temps[1:]))   # [(12, 15), (15, 9), (9, 10)]
```

Because `zip` stops at the shorter input, the last element is never paired
with a missing neighbour. Lists of length 0 or 1 simply produce no pairs.

The index version needs care with the bounds:

```python
for i in range(len(temps) - 1):
    a, b = temps[i], temps[i + 1]
```

## Useful variations

- `all(a <= b for a, b in zip(xs, xs[1:]))` checks that a list is non-decreasing.
- `itertools.pairwise(xs)` (3.10+) yields the same pairs without slicing a copy.
- `max(..., default=0)` handles the case where there are no pairs at all.
- `sum(b > a for a, b in zip(xs, xs[1:]))` counts the number of increases.

## Interview relevance

Neighbour comparisons appear in Best Time to Buy and Sell Stock II, Monotonic
Array, Longest Continuous Increasing Subsequence, and "is this sorted" checks.
The `zip(xs, xs[1:])` idiom removes off-by-one risk at the end of the list and
keeps the solution to a single readable expression.
