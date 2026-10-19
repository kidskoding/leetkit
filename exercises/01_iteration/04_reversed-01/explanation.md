# reversed()

`reversed(seq)` returns an iterator that walks a sequence from the last item to
the first. It does not copy the data.

```python
for name in reversed(events):
    ...

for i in reversed(range(len(events))):
    ...  # len(events) - 1 down to 0
```

Compare with slicing:

```python
for name in events[::-1]:   # builds a full reversed copy first
    ...
```

`events[::-1]` allocates a new list of the same size. `reversed()` is lazy, so
an early `return` costs nothing extra.

## Useful variations

- `reversed(range(n))` gives descending indices without fiddling with
  `range(n - 1, -1, -1)`.
- `next((i for i in reversed(range(len(xs))) if xs[i] == t), -1)` finds the
  last match in one expression.
- `"".join(reversed(s))` reverses a string (though `s[::-1]` is the usual idiom
  when you need the whole result).

## Interview relevance

Use `reversed` for "last occurrence", "most recent", and any right-to-left scan
such as building suffix arrays or processing a stack-like list from the top.
It signals that you know the difference between iterating backwards and
materialising a reversed copy.
