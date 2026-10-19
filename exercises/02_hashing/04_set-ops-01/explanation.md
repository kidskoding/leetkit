# Set operators

Sets support algebra directly through operators:

```python
a = {"ana", "ben", "cal"}
b = {"ben", "dee"}

a & b   # {'ben'}                       intersection
a | b   # {'ana', 'ben', 'cal', 'dee'}  union
a - b   # {'ana', 'cal'}                difference
a ^ b   # {'ana', 'cal', 'dee'}         symmetric difference
```

Each runs in roughly O(len(a) + len(b)). The loop-based equivalent for
intersection is longer and, if `b` is a list, O(n * m):

```python
shared = [x for x in a if x in b]
```

Sets are unordered, so wrap results in `sorted()` whenever the order of the
output matters.

## Useful variations

- Method forms accept any iterable: `a.intersection(some_list)`, `a.union(x, y)`.
- In-place updates: `a &= b`, `a |= b`, `a -= b`.
- Comparisons: `a <= b` (subset), `a.isdisjoint(b)` (no shared items).

## Interview relevance

Intersection of Two Arrays, Find the Difference, "which users follow each other"
and many string problems reduce to one set operation. Using the operators keeps
the solution short enough that you can spend the interview on the follow-ups
(duplicates, sorted output, memory limits) instead of on loops.
