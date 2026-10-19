# Counter & and elements()

`a & b` builds a Counter holding the minimum of each key's counts, and drops
keys whose minimum is zero. `elements()` goes the other way: it yields each
key repeated as many times as its count.

```python
from collections import Counter

common = Counter("bella") & Counter("label")  # Counter({'l': 2, 'b': 1, 'e': 1, 'a': 1})
common &= Counter("roller")                     # Counter({'l': 2, 'e': 1})
sorted(common.elements())                       # ['e', 'l', 'l']
```

Without these operators you would track a dict of minimums by hand, then loop
over it with `for ch, n in mins.items(): result.extend([ch] * n)`.

## Useful variations

- `a | b` keeps the maximum of each count, the multiset union.
- `functools.reduce(operator.and_, map(Counter, words))` folds the whole list in one expression.
- `elements()` skips keys with zero or negative counts, so it is safe after subtraction too.

## Interview relevance

Find Common Characters, Intersection of Two Arrays II, and "items every user
has" questions are all multiset intersections. Reaching for `&` and
`elements()` shows you know Counter is a multiset, not just a frequency dict.
