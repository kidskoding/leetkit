# Counting runs from their starts

Put the numbers in a set, then only start counting at a number whose
predecessor is missing. That number is the smallest value of its run:

```python
num_set = set(nums)
best = 0
for num in num_set:
    if num - 1 not in num_set:
        length = 1
        while num + length in num_set:
            length += 1
        best = max(best, length)
return best
```

The `while` loop looks like it makes this O(n^2), but it only runs from run
starts, and every number belongs to exactly one run. Across the whole loop each
value is stepped over once, so the total work is O(n).

Without the `num - 1 not in num_set` guard, a run of length k is re-counted from
each of its k members, which is quadratic on an input like `range(100000)`.

## Useful variations

- Loop over `num_set`, not `nums`, so duplicates in a huge list do not repeat the start check.
- Sorting first is O(n log n) and also works: `sorted(set(nums))`, then count adjacent gaps of 1.
- To return the run itself, keep the best start too: `list(range(start, start + best))`.

## Interview relevance

Longest Consecutive Sequence is a standard "beat sorting" question. The set
gives O(1) membership, and the "only start at a run's beginning" observation is
the real insight. Be ready to argue the amortised O(n) bound: interviewers often
ask "isn't the inner while loop quadratic?" to see whether you can.
