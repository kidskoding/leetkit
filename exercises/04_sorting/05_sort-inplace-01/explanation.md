# list.sort() versus sorted()

Python has two ways to sort, and they differ in what they touch and return.

```python
nums = [3, 1, 2]
new = sorted(nums)        # new == [1, 2, 3], nums unchanged
nums.sort()               # nums is now [1, 2, 3], returns None
```

Both accept `key=` and `reverse=True`:

```python
sorted(scores, reverse=True)[:k]   # k largest, largest first
```

A classic bug is `result = nums.sort()`, which leaves `result` as `None`.
Another is sorting a function argument in place, which reorders the caller's
list behind their back.

## Useful variations

- `ranked = scores.copy()` then `ranked.sort(reverse=True)` sorts a copy in
  place, the same result as `sorted`.
- `nums.sort()` on a list you own saves the extra O(n) memory of a copy.
- `sorted(text)` works on any iterable (strings, sets, dict keys) and always
  returns a list; `.sort()` exists only on lists.
- `heapq.nlargest(k, scores)` is O(n log k) when k is much smaller than n.

## Interview relevance

Sorting a copy versus sorting in place comes up in Kth Largest Element, Top K
problems and any helper that receives a caller's list. Saying "I'll use
sorted() so I don't mutate the input" or "I'll sort in place to save memory"
shows you understand side effects, and avoiding `x = nums.sort()` avoids a
confusing `None` bug in the middle of a live interview.
