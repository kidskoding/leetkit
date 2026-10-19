# dict(zip(...)) and strict=True

`dict()` accepts any iterable of `(key, value)` pairs, and `zip` produces
exactly that from two parallel lists.

```python
items = ["apple", "bread"]
prices = [3, 5]
dict(zip(items, prices))   # {'apple': 3, 'bread': 5}
```

Plain `zip` stops at the shorter input, so a missing price silently drops an
item. Since Python 3.10, `strict=True` turns that into an error:

```python
dict(zip(["a", "b", "c"], [1, 2], strict=True))
# ValueError: zip() argument 2 is shorter than argument 1
```

The manual version needs its own length check and index loop:

```python
if len(items) != len(prices):
    raise ValueError("length mismatch")
lookup = {}
for i in range(len(items)):
    lookup[items[i]] = prices[i]
```

## Useful variations

- `{k: v for k, v in zip(keys, values) if v is not None}` filters while building.
- `dict(zip(names, range(len(names))))` maps each name to its index (same as
  `{n: i for i, n in enumerate(names)}`).
- Duplicate keys keep the last value; use `setdefault` if the first should win.
- `itertools.zip_longest(a, b, fillvalue=0)` pads instead of truncating or raising.

## Interview relevance

Parallel arrays show up in Evaluate Division, "Employee Importance", sorting
people by heights, and many design questions. `dict(zip(keys, values))` builds
the lookup in one expression, and mentioning `strict=True` shows you think about
malformed input instead of letting `zip` hide it.
