# divmod()

`divmod(a, b)` returns the tuple `(a // b, a % b)`: quotient and remainder.

```python
minutes, secs = divmod(3725, 60)   # (62, 5)
hours, minutes = divmod(minutes, 60)  # (1, 2)
```

It replaces the pair of separate operations:

```python
secs = seconds % 60
minutes = seconds // 60 % 60
hours = seconds // 3600
```

The `divmod` version names each step, and you cannot mismatch the divisor
between the `//` and the `%`.

## Useful variations

- `row, col = divmod(index, num_cols)` maps a flat index to grid coordinates.
- Digit loop: `while n: n, digit = divmod(n, 10)` peels digits from the right.
- Base conversion: `n, r = divmod(n, base)` repeatedly, collecting each `r`.
- `f"{minutes:02d}"` pads to two digits with leading zeros.

## Interview relevance

Add Two Numbers (carry, digit = divmod(total, 10)), Reverse Integer, Excel
column titles, and flattening 2D matrices all need quotient and remainder
together. `divmod` makes that intent obvious in one line.
