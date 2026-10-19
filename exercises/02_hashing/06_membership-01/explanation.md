# Sets for fast membership

`x in container` costs different amounts depending on the container:

- list, tuple, str: O(n), it scans element by element
- set, dict: O(1) on average, it hashes `x` and jumps to the slot

So when you check membership many times against the same collection, convert it
to a set once:

```python
jewel_set = set(jewels)
return sum(1 for ch in stones if ch in jewel_set)
```

Without the conversion, `ch in jewels` rescans the jewels string for every
stone. For small inputs you will not notice; for 10^5 by 10^5 you will time out.

## Useful variations

- `sum(ch in jewel_set for ch in stones)` also works, because `True` counts as 1.
- `frozenset(...)` is an immutable set, useful as a constant or a dict key.
- A dict gives the same O(1) membership when you also need a value per key.

## Interview relevance

Jewels and Stones is an easy warm-up, but the idea behind it is everywhere:
filtering by a banned-word list, checking whether a neighbour is in a word list
in Word Ladder, or testing `num - 1 in nums_set` in Longest Consecutive Sequence.
Saying "I will convert this to a set for O(1) lookups" is exactly the complexity
reasoning interviewers look for.
