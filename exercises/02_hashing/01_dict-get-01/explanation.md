# dict.get()

`d.get(key, default)` returns `d[key]` if the key exists, otherwise `default`.
It never raises `KeyError`.

```python
counts = {}
for ch in s:
    counts[ch] = counts.get(ch, 0) + 1
```

Compare with the branching version:

```python
if ch in counts:
    counts[ch] += 1
else:
    counts[ch] = 1
```

Both are correct, but `get` collapses four lines into one and states the intent
directly: "the count so far, or zero".

## Useful variations

- `d.get(key)` with no default returns `None` for a missing key.
- `Counter(s)` from `collections` builds the same map in one call.
- `defaultdict(int)` lets you write `counts[ch] += 1` with no default at all.
- `max(counts, key=counts.get)` finds the most frequent key.

## Interview relevance

Frequency maps show up in Valid Anagram, Ransom Note, Top K Frequent Elements,
and nearly every sliding-window problem. Knowing `get` (and its cousins
`Counter` and `defaultdict`) lets you write the counting step without thinking,
so your attention stays on the actual algorithm.
