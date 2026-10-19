# Counter.most_common()

`counter.most_common(k)` returns the `k` highest-count entries as a list of
`(element, count)` tuples, most frequent first.

```python
from collections import Counter

Counter("mississippi").most_common(2)  # [('i', 4), ('s', 4)]
```

Elements with equal counts keep the order in which they were first
encountered, which is why `'i'` comes before `'s'` above.

It replaces the manual version:

```python
counts = Counter(nums)
ordered = sorted(counts, key=lambda x: -counts[x])[:k]
```

## Useful variations

- `most_common()` with no argument returns every entry, sorted by count.
- `most_common(1)[0][0]` is the single most frequent element (the mode).
- `most_common()[:-k-1:-1]` gives the `k` least common entries.

## Interview relevance

Top K Frequent Elements, Top K Frequent Words, and "most common word in a
paragraph" all reduce to counting plus ranking. `most_common(k)` gets you a
correct answer quickly; be ready to discuss that it sorts (O(n log n)) and
that a heap (`heapq.nlargest`) or bucket sort can do better for large inputs.
