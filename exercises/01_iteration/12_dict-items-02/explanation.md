# Ranking with sorted(d.items(), key=...)

`d.items()` yields `(key, value)` pairs, and `sorted` accepts a `key` function
that decides what each pair is compared by. Returning a tuple gives a primary
order with a tie-break.

```python
sales = {"tea": 30, "gum": 30, "ink": 50}
sorted(sales.items(), key=lambda kv: (-kv[1], kv[0]))
# [('ink', 50), ('gum', 30), ('tea', 30)]
```

`-kv[1]` puts larger counts first; when counts tie, Python compares the next
tuple field, `kv[0]`, alphabetically. Slicing with `[:n]` keeps the top n and
quietly handles n larger than the list.

The manual version builds the pairs, sorts twice, and trims:

```python
pairs = [(name, sales[name]) for name in sales]
pairs.sort(key=lambda p: p[0])
pairs.sort(key=lambda p: p[1], reverse=True)   # stable sort keeps name order
top = pairs[:n]
```

## Useful variations

- `max(d.items(), key=lambda kv: kv[1])` finds the single best entry.
- `heapq.nsmallest(n, d.items(), key=lambda kv: (-kv[1], kv[0]))` is O(len * log n)
  when n is small.
- `Counter(words).most_common(n)` ranks by count but breaks ties by first
  appearance, not alphabetically.
- Negation only works for numbers; for descending strings, sort twice with
  `reverse=True` and rely on stability.

## Interview relevance

Top K Frequent Words, leaderboards, Sort Characters by Frequency, and Rank Teams
by Votes all finish with "order these counts, highest first, ties by name". A
tuple key with a negated count is the shortest correct way to say that, and
interviewers often probe the tie-break specifically.
