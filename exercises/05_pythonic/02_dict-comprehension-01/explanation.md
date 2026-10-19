# Dict comprehensions

A dict comprehension builds a dictionary from an iterable in one expression.

```python
lengths = {word: len(word) for word in words}
```

It replaces the loop that creates an empty dict and assigns keys one by one:

```python
lengths = {}
for word in words:
    lengths[word] = len(word)
```

If the same key is produced twice, the later value wins, exactly as with
repeated assignment in the loop.

## Useful variations

- `{x: i for i, x in enumerate(nums)}` maps each value to its (last) index.
- `{v: k for k, v in mapping.items()}` inverts a dict.
- `{k: v for k, v in scores.items() if v >= 50}` filters a dict.
- `dict(zip(keys, values))` is the shorter choice when you already have two
  parallel lists.

## Interview relevance

Many solutions start by building a lookup table: value to index for Two Sum,
character to position for a sliding window, id to object for graph problems.
A dict comprehension does that setup in one line and makes the intent obvious.
Remember the "last write wins" rule when duplicates matter, for example when
you need the first index instead of the last.
