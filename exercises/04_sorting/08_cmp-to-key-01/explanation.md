# functools.cmp_to_key

Python 3's `sorted` only takes a `key`, but some orders are easiest to state as
"should `a` come before `b`?". `cmp_to_key` wraps a two-argument comparator so
`sorted` can use it. The comparator returns a negative number when `a` goes
first, a positive number when `b` goes first, and `0` for a tie.

```python
from functools import cmp_to_key

def compare(a, b):
    if a + b > b + a:
        return -1
    if a + b < b + a:
        return 1
    return 0

sorted(["3", "30", "34", "5", "9"], key=cmp_to_key(compare))
# ['9', '5', '34', '3', '30']
```

No per-element key fits Largest Number: `"3"` beats `"30"` but loses to `"34"`,
and only comparing the two concatenations captures that. The comparator must be
consistent (transitive), which `a + b` versus `b + a` is.

## Useful variations

- `(x > y) - (x < y)` is a compact comparator that returns -1, 0 or 1.
- `sorted(words, key=cmp_to_key(cmp), reverse=True)` flips any comparator.
- `list.sort(key=cmp_to_key(cmp))` sorts in place with the same comparator.
- Prefer a plain `key` (often a tuple) whenever one exists; it is faster and simpler.

## Interview relevance

Largest Number, custom alien-dictionary orders and "sort by a rule that
depends on both items" questions need a comparator. Knowing that `cmp=` is
gone in Python 3 and reaching for `cmp_to_key` immediately avoids a fumble,
and the "0" edge case shows you test the boundaries of your output.
