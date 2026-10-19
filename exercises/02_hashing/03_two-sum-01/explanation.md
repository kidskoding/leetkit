# Complement lookup with a dict

Store each value's index as you go, and before storing, ask whether the value
you need has already appeared:

```python
seen = {}
for i, num in enumerate(nums):
    complement = target - num
    if complement in seen:
        return [seen[complement], i]
    seen[num] = i
```

The brute-force version checks every pair:

```python
for i in range(len(nums)):
    for j in range(i + 1, len(nums)):
        if nums[i] + nums[j] == target:
            return [i, j]
```

The dict turns the inner loop into a single O(1) lookup, so the whole thing
runs in O(n).

## Useful variations

- Look up before inserting: inserting first would let `[3, 2, 4], 6` pair index 0 with itself.
- Returning `[seen[complement], i]` keeps the indices in increasing order for free.
- The same shape solves "count pairs with sum k": add `seen.get(complement, 0)` to a total.

## Interview relevance

Two Sum is often the first question in a screen, and the complement-in-a-dict
idea generalises to Subarray Sum Equals K (prefix sums), 4Sum II and pair-counting
problems. Explaining why you check before inserting is a quick way to show you
understand the edge cases.
