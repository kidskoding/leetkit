# Sets for duplicate detection

A `set` stores unique, hashable values with average O(1) `add` and `in`.
Tracking "what have I seen so far" is the classic use:

```python
seen = set()
for x in items:
    if x in seen:
        return True
    seen.add(x)
return False
```

Checking `x in some_list` instead would scan the whole list each time,
turning the loop into O(n^2).

There is also a one-liner when you do not need an early exit:

```python
return len(set(items)) != len(items)
```

It always processes the full input, while the loop stops at the first repeat.

## Useful variations

- `seen = set()` creates an empty set; `{}` is an empty dict, not a set.
- `set(items)` deduplicates a list in one call (order is not preserved).
- `seen.add(x)` returns `None`, so do not write `seen = seen.add(x)`.

## Interview relevance

Contains Duplicate is a frequent opener, and the seen-set loop is the core of
Happy Number, Linked List Cycle (by node), graph BFS/DFS visited tracking and
Longest Substring Without Repeating Characters. Mentioning the O(n) time versus
O(n) extra space trade-off against sorting shows you know why the set helps.
