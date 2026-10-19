# Monotonic deque

A monotonic deque keeps indices whose values only decrease from front to
back. When a new value arrives, every smaller value behind it is useless (it
leaves the window earlier and is never larger), so pop those from the right.
When the front index slides out of the window, pop it from the left.

```python
from collections import deque

window = deque()
for i, x in enumerate(nums):
    while window and nums[window[-1]] <= x:
        window.pop()
    window.append(i)
    if window[0] <= i - k:
        window.popleft()
    if i >= k - 1:
        maxima.append(nums[window[0]])
```

The clunky version is `[max(nums[i:i + k]) for i in range(len(nums) - k + 1)]`,
which is O(nk). Each index enters and leaves the deque at most once, so the
deque version is O(n).

## Useful variations

- Flip the comparison to `>=` to track window minimums instead.
- Store indices, not values, so you can tell when the front has left the window.
- Combine a max deque and a min deque to find the longest window whose max minus min stays within a limit.

## Interview relevance

Sliding Window Maximum is a standard hard problem, and the same structure
powers Longest Continuous Subarray With Absolute Diff Limit and Shortest
Subarray with Sum at Least K. A list cannot pop from the front in O(1), which
is why `deque` with `pop` and `popleft` is the right tool.
