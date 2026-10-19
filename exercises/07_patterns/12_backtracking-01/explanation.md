# Backtracking with append and pop

Backtracking builds candidates one choice at a time in a single shared list.
Append a choice, recurse, then pop it to undo the choice before trying the
next one. A nested function can read and mutate `path` and `loads` from the
enclosing scope without passing them around.

```python
weights = sorted(weights)
loads, path = [], []

def backtrack(start, remaining):
    if remaining == 0:
        loads.append(path[:])
        return
    for i in range(start, len(weights)):
        if weights[i] > remaining:
            break
        path.append(weights[i])
        backtrack(i, remaining - weights[i])
        path.pop()

backtrack(0, target)
```

`path[:]` records a copy; appending `path` itself would store a reference that
later pops empty out. Starting the loop at `start` stops the same multiset
appearing in two orders, and passing `i` (not `i + 1`) lets an item be reused.
Because the weights are sorted, `break` prunes every heavier choice at once.
Building a new list with `path + [w]` on each call also works but copies on
every step.

## Useful variations

- Pass `i + 1` instead of `i` when each item may be used at most once (Combination Sum II, Subsets).
- Record `path[:]` at the top of every call, not only at the target, to generate all subsets.
- Skip `if i > start and nums[i] == nums[i - 1]: continue` to drop duplicates in sorted input.
- `itertools.combinations` covers fixed-size choices without reuse when no pruning is needed.

## Interview relevance

Combination Sum, Subsets, Permutations, Letter Combinations, Word Search, and
N-Queens all follow choose, explore, un-choose. Interviewers look for the copy
when recording, a clean base case, and pruning that keeps the search small.
