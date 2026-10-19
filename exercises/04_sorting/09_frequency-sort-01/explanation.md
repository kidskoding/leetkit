# Counter plus a tuple sort key

`Counter` turns a list into a word-to-count mapping in one call. Sorting it
iterates over the distinct words, and a tuple key orders them by several
criteria at once: tuples compare element by element, so the second item only
matters when the first ties.

```python
from collections import Counter

words = ["kiwi", "banana", "kiwi", "banana", "date"]
counts = Counter(words)
sorted(counts, key=lambda w: (-counts[w], w))
# ['banana', 'kiwi', 'date']
```

Negating the count makes larger counts sort first while the word still sorts
ascending. `reverse=True` would flip both, giving reverse-alphabetical ties.

## Useful variations

- `sorted(counts.items(), key=lambda kv: (-kv[1], kv[0]))` keeps the counts as `(word, count)` pairs.
- `counts.most_common(k)` gives the top `k` by count, but ties keep first-seen order, not alphabetical.
- `heapq.nsmallest(k, counts, key=lambda w: (-counts[w], w))` picks the top `k` without a full sort.
- `Counter(s)` on a string counts characters, ready for Sort Characters by Frequency.

## Interview relevance

Top K Frequent Words asks for exactly this order, and similar "most common,
then alphabetical" rules show up in leaderboards and log summaries. Writing
`Counter` plus a `(-count, word)` key in two lines is the expected Python
answer, and it shows you know why `most_common` alone does not handle ties.
