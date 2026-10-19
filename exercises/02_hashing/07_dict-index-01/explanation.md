# Dict insertion order

Since Python 3.7, a `dict` iterates its keys in the order they were first
inserted. Updating an existing key does not move it. That turns a dict into an
"ordered record of what I have seen":

```python
positions = {}
for i, ch in enumerate(s):
    positions[ch] = -1 if ch in positions else i

for i in positions.values():
    if i != -1:
        return i
return -1
```

The second loop walks characters in order of first appearance, so the first
non-repeated one it finds is the answer. Compare with the quadratic version:

```python
for i, ch in enumerate(s):
    if s.count(ch) == 1:
        return i
```

`s.count(ch)` scans the whole string each time, giving O(n^2).

## Useful variations

- Count first with `Counter(s)`, then `for i, ch in enumerate(s): if counts[ch] == 1`.
- `next((i for i in positions.values() if i != -1), -1)` is the one-expression form of the search.
- `list(d)` gives keys in insertion order; `dict.fromkeys(items)` dedupes while keeping order.
- `set` does not preserve order, so do not use it when the answer depends on "first".

## Interview relevance

"First unique", "first repeated" and "deduplicate but keep order" questions all
lean on ordered hashing. Knowing that a plain dict is ordered (and a set is
not) lets you skip `OrderedDict` and avoid a second pass over the input,
which interviewers often ask about as a follow-up for streams of characters.
