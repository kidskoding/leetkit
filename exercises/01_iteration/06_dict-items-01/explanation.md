# dict.items()

`d.items()` yields `(key, value)` pairs, which you unpack directly in the loop.

```python
for category, amount in spending.items():
    ...
```

Compare with looping over keys and indexing back in:

```python
for category in spending:
    amount = spending[category]
```

The `.items()` version does one lookup instead of two and makes it obvious that
both pieces are used. Since Python 3.7, dicts preserve insertion order, so the
iteration order is the order keys were first added.

## Useful variations

- `d.keys()` and `d.values()` when you only need one side.
- `{k: v for k, v in d.items() if v > 0}` filters a dict into a new dict.
- `max(d.items(), key=lambda kv: kv[1])` finds the pair with the largest value.
- `sorted(d.items(), key=lambda kv: -kv[1])` orders pairs by value, descending.

## Interview relevance

Frequency-count problems (Top K Frequent Elements, First Unique Character,
Majority Element) and grouping problems (Group Anagrams) all finish by scanning
a dict. Writing `for key, value in d.items()` is the expected idiom, and
iterating while mutating the same dict is a classic bug to avoid.
