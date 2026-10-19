# defaultdict(list)

`defaultdict(factory)` is a dict that calls `factory()` to create a value the
first time a missing key is accessed.

```python
from collections import defaultdict

groups = defaultdict(list)
for word in words:
    groups["".join(sorted(word))].append(word)
```

Compare with a plain dict:

```python
groups = {}
for word in words:
    key = "".join(sorted(word))
    if key not in groups:
        groups[key] = []
    groups[key].append(word)
```

Dicts (including `defaultdict`) keep insertion order, so `list(groups.values())`
returns groups in the order their first member was seen.

## Useful variations

- `tuple(sorted(word))` also works as a key; lists cannot be keys because they are not hashable.
- A 26-slot letter-count tuple as the key avoids sorting each word (O(n) per word).
- `groups.setdefault(key, []).append(word)` does the same with a plain dict.

## Interview relevance

Group Anagrams, Group Shifted Strings, and "bucket records by attribute" all
follow the same shape: compute a hashable signature, append to its bucket.
Writing it with `defaultdict(list)` keeps the solution to three lines and lets
you focus the discussion on choosing the key.
