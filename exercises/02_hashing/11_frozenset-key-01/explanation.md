# frozenset as a dict key

Dict keys must be hashable, and a hash must never change, so mutable containers
like `set` and `list` are rejected. `frozenset` is the immutable twin of `set`:
same membership rules, same "order and repeats do not matter" equality, but
hashable, so it can be a key:

```python
groups = {}
for word in words:
    groups.setdefault(frozenset(word), []).append(word)
return list(groups.values())
```

`frozenset("loop") == frozenset("lop")` is True, so both words land in the same
bucket. Trying the same with a plain set fails immediately:

```python
groups[set(word)] = []   # TypeError: unhashable type: 'set'
```

The workaround people reach for, `"".join(sorted(set(word)))`, works but builds
a sorted string just to get a canonical key; the frozenset says what you mean.

## Useful variations

- Count trades between two people in either direction with a `Counter`: `pairs[frozenset((a, b))] += 1`.
- A tuple key, `tuple(sorted(word))`, is the right choice when repeats DO matter, as in anagram grouping.
- Sets of frozensets work too: `seen.add(frozenset(state))` dedupes unordered states in a search.
- Count Pairs of Similar Strings is this grouping plus `n * (n - 1) // 2` per group.

## Interview relevance

Unordered keys come up whenever direction or order is irrelevant: mutual
friendships, undirected edges, "same characters" string problems, and
memoising over a set of chosen items. Knowing that `set` is unhashable and
`frozenset` is the fix saves a confusing `TypeError` mid-interview.
