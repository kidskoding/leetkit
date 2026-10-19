# Bottom-up DP with a list and math.inf

A bottom-up DP keeps one list indexed by subproblem. Here `best[a]` is the
fewest coins that make amount `a`. `math.inf` marks "not reachable yet": it
compares correctly with `min`, and `math.inf + 1` is still `math.inf`, so no
special cases are needed.

```python
import math

best = [0] + [math.inf] * amount
for a in range(1, amount + 1):
    for coin in coins:
        if coin <= a:
            best[a] = min(best[a], best[a - coin] + 1)
answer = -1 if best[amount] == math.inf else best[amount]
```

The same recurrence top-down reads almost like the maths:

```python
from functools import lru_cache

@lru_cache(maxsize=None)
def fewest(a):
    if a == 0:
        return 0
    return min((fewest(a - c) + 1 for c in coins if c <= a), default=math.inf)
```

Without the cache the recursion is exponential. The decorator memoises each
`a`, giving the same O(amount * len(coins)) work as the table.

## Useful variations

- `functools.cache` (Python 3.9+) is shorthand for `lru_cache(maxsize=None)`.
- `min(..., default=math.inf)` handles an empty generator when no coin fits.
- Count the number of ways instead with `ways[a] += ways[a - coin]`, looping coins in the outer loop.
- Top-down recursion can hit Python's ~1000-frame limit on large amounts; the bottom-up list cannot.

## Interview relevance

Coin Change, Coin Change II, Perfect Squares, Word Break, and Minimum Cost for
Tickets share this one-dimensional table. Interviewers expect the recurrence,
a sentinel for impossible states, a counterexample to greedy, and awareness of
the recursion-depth trade-off between memoised recursion and a table.
