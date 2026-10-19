# enumerate(iterable, start=1)

`enumerate` takes an optional `start` argument that sets the first number it
hands out. The values are unchanged; only the counter shifts.

```python
for n, line in enumerate(lines, start=1):
    print(f"{n}. {line}")
```

Without `start`, the 1-based number has to be patched in by hand:

```python
for i, line in enumerate(lines):
    print(f"{i + 1}. {line}")
```

That `+ 1` is easy to forget in one place and duplicate in another, especially
when the number is used more than once in the loop body.

## Useful variations

- `enumerate(lines, 1)` is the positional form; `start=1` reads more clearly.
- `{name: place for place, name in enumerate(finishers, start=1)}` maps each
  runner to a 1-based rank.
- `for case, data in enumerate(tests, start=1): print(f"Case #{case}: ...")`
  matches contest-style output.
- `start` can be any int, for example `enumerate(rows, start=2)` when row 1 is a
  header.

## Interview relevance

Ranks, line numbers, "return the 1-indexed position" answers (Two Sum II), and
contest output headers all want counting from 1. Using `start=1` keeps the
indexing convention in one visible place instead of scattering `i + 1` through
the code, which is a classic source of off-by-one bugs under time pressure.
