# Two pointers on a sorted array

Two indices start at opposite ends and walk toward each other. Sorted order
tells you which one to move: a sum that is too small needs a bigger left value,
a sum that is too big needs a smaller right value.

```python
lo, hi = 0, len(prices) - 1
while lo < hi:
    total = prices[lo] + prices[hi]
    if total == target:
        return (lo, hi)
    if total < target:
        lo += 1
    else:
        hi -= 1
```

Python helps in small ways: tuple assignment sets both pointers in one line,
`while lo < hi` reads exactly like the invariant, and returning a tuple gives the
caller something to unpack with `i, j = pair_with_sum(...)`.

## Useful variations

- `lo, hi = lo + 1, hi - 1` moves both pointers at once (palindromes, reversal).
- `s = [c.lower() for c in s if c.isalnum()]` cleans input before a palindrome check.
- `cleaned == cleaned[::-1]` is the one-liner check when O(n) extra space is fine.
- `xs[lo], xs[hi] = xs[hi], xs[lo]` swaps in place for reverse-style problems.

## Interview relevance

Two Sum II, 3Sum, container with most water, valid palindrome, and merging
sorted arrays are all two-pointer problems. The pattern turns an O(n^2) pair
search into O(n) with O(1) space, and interviewers expect you to justify why
moving one pointer never discards the answer.
