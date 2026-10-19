# itertools.combinations()

`combinations(iterable, r)` yields every group of `r` items, taken by position,
with no repeats and in input order.

```python
from itertools import combinations

list(combinations([3, 5, 2], 2))  # [(3, 5), (3, 2), (5, 2)]
```

It replaces the nested index loop:

```python
for i in range(len(prices)):
    for j in range(i + 1, len(prices)):
        a, b = prices[i], prices[j]
```

Items are chosen by position, not value, so `[4, 4, 4]` gives three `(4, 4)`
pairs. The output order is lexicographic by index: all pairs starting at
index 0, then index 1, and so on.

## Useful variations

- `combinations(range(n), 2)` yields index pairs `(i, j)` when you need positions.
- `permutations(items, 2)` yields ordered pairs, so `(a, b)` and `(b, a)` both appear.
- `combinations_with_replacement(items, 2)` also pairs each item with itself.
- `math.comb(n, r)` counts the combinations without generating them.

## Interview relevance

Saying "the brute force checks every pair" and writing it with `combinations`
takes seconds, leaving time for the optimised hash-map or two-pointer version.
It also generates candidate subsets in small-input problems and test harnesses.
