# collections.Counter

`Counter(iterable)` counts how many times each element appears and stores the
result in a dict subclass.

```python
from collections import Counter

Counter("banana")  # Counter({'a': 3, 'n': 2, 'b': 1})
```

Two Counters are equal when they hold the same counts, so a multiset
comparison is one line:

```python
Counter(s) == Counter(t)
```

The hand-rolled version needs a dict, a loop with `counts.get(ch, 0) + 1`, and
a second loop to compare. Counter also returns `0` for missing keys instead of
raising `KeyError`.

## Useful variations

- `Counter(words)` works on any iterable, not just strings.
- `counts[ch] += 1` works even when `ch` has never been seen.
- `sorted(s) == sorted(t)` is an O(n log n) alternative that needs no import.

## Interview relevance

Frequency counting is everywhere: anagrams, permutations in a sliding window,
"first unique character", majority element. Interviewers expect you to reach
for `Counter` rather than rebuild a counting dict by hand, and to know that
`==` compares the counts directly.
