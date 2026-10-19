# Tuple keys for multi-key sorting

Tuples compare element by element: the first items decide, and only on a tie
does Python look at the second. So a key that returns a tuple sorts by several
criteria at once.

```python
scores = [("ana", 90), ("bo", 95), ("cy", 90)]
sorted(scores, key=lambda p: (-p[1], p[0]))
# [('bo', 95), ('ana', 90), ('cy', 90)]
```

Negating the score makes the highest score the smallest key, so it sorts
first, while the name still sorts in normal ascending order.

The tempting alternative is wrong:

```python
sorted(scores, key=lambda p: (p[1], p[0]), reverse=True)
```

`reverse=True` flips every field, so tied players come out Z to A.

## Useful variations

- `sorted(counts.items(), key=lambda kv: (-kv[1], kv[0]))` is the Top K
  Frequent Words ordering.
- For a string field that must be descending, sort twice using stability:
  first by the secondary key, then by the primary key.
- `key=lambda iv: (iv[0], -iv[1])` orders intervals by start, longest first on
  ties.

## Interview relevance

Problems regularly specify a primary order and a tie-breaker: frequency then
value, score then name, start then end. A tuple key with a negated numeric
field handles all of these in one readable call, and explaining why
`reverse=True` would break the tie-breaker shows real understanding.
