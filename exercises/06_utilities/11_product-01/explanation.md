# itertools.product()

`product(*iterables, repeat=1)` yields the cartesian product: every way to pick
one item from each input. With `repeat=n` it picks from the same input `n`
times, which is exactly "every string of length n".

```python
from itertools import product

["".join(t) for t in product("ab", repeat=2)]  # ['aa', 'ab', 'ba', 'bb']
list(product(["red", "blue"], ["jeans", "shorts"]))
# [('red', 'jeans'), ('red', 'shorts'), ('blue', 'jeans'), ('blue', 'shorts')]
```

The order is that of nested loops with the first input outermost, so the
first position changes slowest. It replaces loops you would otherwise have to
nest by hand, which is impossible when `n` is only known at runtime:

```python
for a in alphabet:
    for b in alphabet:
        codes.append(a + b)
```

`product(..., repeat=0)` yields one empty tuple, and any empty input makes the
whole product empty.

## Useful variations

- `product(range(rows), range(cols))` walks every `(r, c)` cell of a grid in one loop.
- `product(*(keypad[d] for d in digits))` solves Letter Combinations of a Phone Number.
- `product([0, 1], repeat=n)` enumerates every subset as a bitmask of choices.
- `len(alphabet) ** n` counts the results without generating them.

## Interview relevance

Letter Combinations of a Phone Number, generating all bit strings, and
brute-force checks over every assignment of choices are all cartesian products.
Writing them with `product` is short, clearly correct, and frees time for the
optimised solution.
