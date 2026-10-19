# itertools.permutations()

`permutations(iterable, r=None)` yields every ordering of the items, as tuples.
With `r` it yields every ordered selection of `r` items.

```python
from itertools import permutations

["".join(p) for p in permutations("cat")]
# ['cat', 'cta', 'act', 'atc', 'tca', 'tac']
```

Items are picked by position, not value, so a word with repeated letters
produces repeated tuples. Collect into a set to keep one of each, then sort for
a stable order:

```python
sorted({"".join(p) for p in permutations("aab")})  # ['aab', 'aba', 'baa']
```

This replaces a hand-written backtracking function with a `used` array and a
path list. Remember the size: n items give n! orderings, so this is only
practical for small inputs (about 10 items or fewer).

## Useful variations

- `permutations(items, 2)` yields ordered pairs, so both `(a, b)` and `(b, a)` appear.
- `int(min("".join(p) for p in permutations(digits)))` finds the smallest number made from some digits.
- `next(permutations(items))` is the input order itself; permutations come out in input-position order.
- `math.perm(n, r)` counts the orderings without generating them.

## Interview relevance

Permutations, Permutations II, and "try every order" brute forces (small
travelling-salesman checks, next-arrangement puzzles) all start here. Using
`permutations` for the baseline, and spotting the duplicate trap, shows fluency
before you move on to the optimised backtracking or math solution.
