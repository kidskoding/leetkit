# itertools.groupby()

`groupby(iterable, key=None)` yields `(key, group)` pairs, one for each run of
consecutive items with the same key. `group` is an iterator over that run.

```python
from itertools import groupby

[(k, len(list(g))) for k, g in groupby(["INFO", "INFO", "WARN", "INFO"])]
# [('INFO', 2), ('WARN', 1), ('INFO', 1)]
```

It replaces the manual "previous item plus counter" loop, with its easy-to-miss
final flush after the loop ends.

The important detail: `groupby` only groups adjacent items. INFO appears twice
above because a WARN separates the runs. To group every equal item together,
sort by the same key first:

```python
for k, g in groupby(sorted(levels)):
    ...
```

Each `group` shares the underlying iterator, so consume it (for example with
`list(g)`) before moving to the next key.

## Useful variations

- `"".join(f"{k}{len(list(g))}" for k, g in groupby(s))` is run-length encoding of a string.
- `groupby(sorted(words, key=len), key=len)` buckets words by length.
- `max(len(list(g)) for k, g in groupby(bits) if k == 1)` finds the longest run of 1s.
- For counts across the whole input regardless of order, use `Counter` instead.

## Interview relevance

String Compression, Count and Say, Max Consecutive Ones, and log summarising
are all run-length problems. `groupby` writes them in one comprehension, and
explaining why it needs sorted input to group globally shows you understand
the tool rather than memorised it.
