# Sliding window with a last-seen dict

A sliding window keeps a range `[start, i]` that always satisfies the rule (here:
no repeated characters). The right edge advances every step; the left edge only
moves when the rule breaks.

```python
last_seen = {}
start = best = 0
for i, ch in enumerate(s):
    if last_seen.get(ch, -1) >= start:
        start = last_seen[ch] + 1
    last_seen[ch] = i
    best = max(best, i - start + 1)
```

`enumerate` gives the right edge and the character together, and
`dict.get(ch, -1)` handles "never seen" without an `if ch in last_seen` branch.
The `>= start` check matters: an old index to the left of the window (as in
`"abba"`) must not drag `start` backwards.

## Useful variations

- `Counter` or `defaultdict(int)` tracks counts when the rule is "at most k distinct".
- `while counts[ch] > 1:` shrinks the window one step at a time when you cannot jump.
- `max(best, i - start + 1)` keeps a running best without storing every window.
- Fixed-size windows can use `sum(nums[:k])` once and then add/subtract the edges.

## Interview relevance

Longest substring without repeats, longest repeating character replacement,
minimum window substring, and permutation-in-string are all sliding windows.
They turn O(n^2) substring checks into one O(n) pass, and the dict bookkeeping
is where most bugs hide.
