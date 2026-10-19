# Tuple unpacking and multiple assignment

Python can assign several names in one statement. The right-hand side is
evaluated completely first, then unpacked into the targets left to right.

```python
a, b = 0, 1
a, b = b, a + b     # uses the OLD a and b on the right
```

Without unpacking you need a temporary variable:

```python
nxt = a + b
a = b
b = nxt
```

The same rule gives the one-line swap:

```python
nums[i], nums[j] = nums[j], nums[i]
```

## Useful variations

- `lo, hi = 0, len(nums) - 1` initialises two pointers together.
- `q, r = divmod(n, 10)` unpacks a function's tuple result.
- `for name, score in pairs:` unpacks each tuple inside a loop.
- Unpacking must match in length: `a, b = [1, 2, 3]` raises `ValueError`.

## Interview relevance

Rolling DP (Fibonacci, Climbing Stairs, House Robber), two-pointer setups, and
in-place swaps (reversing an array, partitioning in quicksort, cyclic sort) all
lean on multiple assignment. It removes temporaries and the bugs that come from
updating one variable before the other was read.
