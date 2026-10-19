# Prefix sums with a Counter

If `total` is the sum of the first `j` values, a run ending at `j` sums to `k`
exactly when some earlier prefix sum equals `total - k`. So instead of trying
every start, count how many earlier prefixes had that value.

```python
from collections import Counter

seen = Counter({0: 1})
total = count = 0
for x in changes:
    total += x
    count += seen[total - k]
    seen[total] += 1
```

`Counter` returns 0 for missing keys without inserting them, so the lookup needs
no `if` and no `.get`. Seeding `{0: 1}` stands for the empty prefix, which is
what lets runs starting at index 0 match. Look up before you record `total`,
otherwise `k == 0` would count an empty run.

## Useful variations

- `defaultdict(int)` works too, but `seen[total - k]` inserts a 0 entry on lookup.
- `itertools.accumulate(nums)` builds the full prefix list when you need random access.
- `seen[total % k]` counts subarrays whose sum is divisible by `k`.
- `first_index.setdefault(total, j)` finds the longest run instead of counting them.

## Interview relevance

Subarray sum equals k, continuous subarray sum, contiguous array (equal 0s and
1s), and range-sum queries all rely on prefix sums. Negative numbers break the
sliding-window approach, which is exactly why interviewers ask this one.
