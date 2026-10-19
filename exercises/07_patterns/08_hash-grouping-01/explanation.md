# Grouping with defaultdict(list) and a computed key

The grouping pattern: compute a canonical key for each item so that items that
"belong together" share a key, then bucket by that key.

```python
from collections import defaultdict

groups = defaultdict(list)
for w in words:
    key = tuple((ord(b) - ord(a)) % 26 for a, b in zip(w, w[1:]))
    groups[key].append(w)
return list(groups.values())
```

Without `defaultdict` you write `if key not in groups: groups[key] = []` on every
item. The key must be hashable, so build it as a `tuple` (or `str`), never a
list. `zip(w, w[1:])` walks neighbouring pairs without indexing, and `% 26`
makes the wrap from `"z"` to `"a"` come out non-negative. Dicts keep insertion
order, so `groups.values()` lists groups by first appearance.

## Useful variations

- `"".join(sorted(w))` or `tuple(sorted(w))` is the key for grouping anagrams.
- `email.rsplit("@", 1)[1].lower()` is the key for grouping users by domain.
- `groups.setdefault(key, []).append(w)` does the same with a plain dict.
- `itertools.groupby` only groups adjacent items, so sort by the key first.

## Interview relevance

Group anagrams, group shifted strings, find duplicate files, and bucketing logs
by user all use this pattern. Interviewers look for a correct canonical key
and a hashable type, and often ask how the output is ordered.
