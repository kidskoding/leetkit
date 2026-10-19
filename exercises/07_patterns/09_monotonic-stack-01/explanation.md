# Monotonic stack with a list of indices

A plain Python list is a stack: `append` pushes and `pop()` removes the last
item, both O(1). For "next greater element" questions, keep the indices whose
answer is still unknown. Their temperatures never increase from bottom to top,
so a warmer day pops every colder index it beats.

```python
answer = [0] * len(temps)
stack = []
for i, t in enumerate(temps):
    while stack and temps[stack[-1]] < t:
        j = stack.pop()
        answer[j] = i - j
    stack.append(i)
```

The clunky version scans forward from every day with a nested loop, which is
O(n^2) on a long cold spell. Here each index is pushed once and popped at most
once, so the whole pass is O(n). `stack[-1]` peeks at the top without removing
it, `enumerate` hands you the index you need for the distance, and
`[0] * len(temps)` pre-fills the "never warmer" answer.

## Useful variations

- Use `<=` instead of `<` in the `while` to treat equal values as "greater or equal".
- Store `answer[j] = t` instead of `i - j` to return the next greater value itself.
- Iterate `range(len(nums) * 2)` with `i % len(nums)` for a circular array.
- Flip the comparison (`>`) to find the next smaller element, as in Largest Rectangle in Histogram.

## Interview relevance

Daily Temperatures, Next Greater Element I and II, Online Stock Span, and
Largest Rectangle in Histogram all reduce to this loop. Interviewers expect
you to explain why it is O(n) despite the inner `while`, and storing indices
rather than values is the detail that makes distances easy.
