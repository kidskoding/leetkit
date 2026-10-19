# itertools.accumulate()

`accumulate(iterable)` yields running totals. With `initial=0` (Python 3.8+) it
also yields a leading 0, which is exactly the prefix-sum array you want.

```python
from itertools import accumulate

prefix = list(accumulate([2, 4, 1, 3], initial=0))  # [0, 2, 6, 7, 10]
total = prefix[right + 1] - prefix[left]
```

It replaces the manual running-total loop:

```python
prefix = [0]
for num in nums:
    prefix.append(prefix[-1] + num)
```

## Useful variations

- `accumulate(nums, max)` gives a running maximum; `accumulate(nums, min)` a running minimum.
- `accumulate(nums, operator.mul)` gives running products, useful for Product of Array Except Self.
- Without `initial`, `list(accumulate(nums))` has the same length as `nums` and no leading 0.

## Interview relevance

Range Sum Query, Subarray Sum Equals K, Contiguous Array, and 2D prefix-sum
problems all start from a prefix array. The leading 0 from `initial=0` removes
the `if left == 0` special case that causes off-by-one bugs.
