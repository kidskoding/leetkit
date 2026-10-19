# f-strings with format specs

An f-string evaluates expressions inside `{}`. Anything after a colon is a
format spec that controls alignment, width and precision:

```python
total = qty * unit_price
return f"{item:<10} {qty:>3} {total:>8.2f}"
```

- `<10` left-aligns in a field 10 characters wide.
- `>3` right-aligns in width 3.
- `>8.2f` right-aligns in width 8 and rounds to two decimal places.

The hand-rolled version needs `str()`, `round()` and space counting, and
`str(round(1.5, 2))` gives `"1.5"`, not `"1.50"`:

```python
price = str(round(total, 2))
return item + " " * (10 - len(item)) + " " + ...
```

Widths are minimums: a longer value is printed in full and pushes the rest of
the line over.

## Useful variations

- `f"{name:^12}"` centers; `f"{n:*>6}"` pads with `*` instead of spaces.
- `f"{n:05d}"` zero-pads integers, and `f"{n:,}"` adds thousands separators.
- `f"{ratio:.1%}"` formats `0.256` as `25.6%`.
- `f"{x=}"` prints both the expression and its value, handy for debugging.

## Interview relevance

Output formatting shows up in design and simulation questions (receipts,
leaderboards, logs) and whenever you print intermediate state while debugging
live. Knowing format specs means one readable line instead of fiddly string
arithmetic, and it avoids the classic `1.5` versus `1.50` bug.
