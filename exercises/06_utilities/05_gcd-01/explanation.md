# math.gcd() with functools.reduce()

`math.gcd(a, b)` returns the greatest common divisor of two integers.

```python
import math

math.gcd(12, 18)  # 6
math.gcd(0, 5)    # 5
```

To combine it across a list, fold it with `functools.reduce`:

```python
from functools import reduce

reduce(math.gcd, [12, 18, 30], 0)  # gcd(gcd(gcd(0, 12), 18), 30) == 6
```

Starting from `0` works because `gcd(0, x) == x`, and it makes an empty list
return `0` instead of raising `TypeError`.

## Useful variations

- `math.gcd(*counts)` accepts any number of arguments in Python 3.9+.
- `math.lcm(a, b)` (3.9+) gives the least common multiple; `a * b // math.gcd(a, b)` works everywhere.
- Simplify a fraction: `g = math.gcd(num, den)`, then `(num // g, den // g)`.

## Interview relevance

GCD is the hidden step in Greatest Common Divisor of Strings, X of a Kind in a
Deck of Cards, fraction arithmetic, and "can these jugs measure x litres".
Reaching for `math.gcd` instead of hand-writing Euclid's algorithm keeps the
focus on the real problem.
