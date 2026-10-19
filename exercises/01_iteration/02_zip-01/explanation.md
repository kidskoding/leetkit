# zip()

`zip(a, b)` yields tuples `(a[0], b[0]), (a[1], b[1]), ...` and stops as soon
as the shortest input runs out.

```python
for given, correct in zip(answers, key):
    ...
```

It replaces the index-based loop:

```python
for i in range(min(len(answers), len(key))):
    given, correct = answers[i], key[i]
```

The `zip` version needs no length arithmetic, works with any iterables, and
reads as "these two things go together".

## Useful variations

- `sum(x == y for x, y in zip(a, b))` counts matches, since `True` is `1`.
- `dict(zip(names, scores))` builds a mapping from two parallel lists.
- `zip(a, b, c)` walks three (or more) sequences at once.
- `zip(a, b, strict=True)` (3.10+) raises if the lengths differ instead of
  silently truncating.

## Interview relevance

Any time two sequences are compared position by position, `zip` is the idiom:
Hamming distance, Bulls and Cows, checking if two strings differ by one edit,
pairing names with values. Remember the truncation rule: if unequal lengths
matter to the answer, handle them explicitly or use `strict=True`.
