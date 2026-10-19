# functools.lru_cache() and functools.cache()

`@lru_cache(maxsize=None)` (or its shorthand `@cache`, Python 3.9+) wraps a
function so each distinct argument is computed once; later calls return the
stored result.

```python
from functools import cache

@cache
def climb_stairs(n):
    if n <= 1:
        return 1
    return climb_stairs(n - 1) + climb_stairs(n - 2)
```

Without the decorator the recursion branches twice per call and recomputes
the same subproblems over and over: `climb_stairs(30)` makes about 2.7 million
calls, and `climb_stairs(90)` would make around 10^19, which never finishes.
With the cache each `n` is solved once, so `climb_stairs(90)` takes 91 calls
and is instant. It replaces the hand-written memo dict:

```python
memo = {}
def climb(n):
    if n in memo:
        return memo[n]
    ...
```

Arguments must be hashable (ints, strings, tuples), so convert lists to tuples
before passing them in.

## Useful variations

- `@lru_cache(maxsize=128)` bounds memory by evicting the least recently used results.
- Define the cached helper inside the main function (`@cache def dp(i, j): ...`) so the cache resets per call.
- `climb_stairs.cache_clear()` empties the cache; `climb_stairs.cache_info()` shows hits and misses.
- Very deep recursions can still hit the recursion limit; switch to a bottom-up loop if n is in the thousands.

## Interview relevance

Climbing Stairs, Unique Paths, Coin Change, Word Break, and Longest Common
Subsequence are commonly solved top-down. Writing the recurrence and adding
one decorator gets you a correct memoized solution fast, and it is easy to
explain as the stepping stone to a bottom-up table.
