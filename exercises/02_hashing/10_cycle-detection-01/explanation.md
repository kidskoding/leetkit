# Cycle detection with a seen set

When a process is deterministic (the next state depends only on the current
one), revisiting any state means it will loop forever. A set of seen states
turns "does this ever stop?" into an O(1) check per step:

```python
seen = set()
while n != 1 and n not in seen:
    seen.add(n)
    n = sum(int(d) ** 2 for d in str(n))
return n == 1
```

The loop ends for one of two reasons, and `n == 1` tells you which. Without the
set, an unhappy number like 2 would spin forever, and a step limit such as
"give up after 1000 iterations" is a guess, not a proof.

`str(n)` is the quickest way to get the digits in Python. The arithmetic
version is `while n: n, d = divmod(n, 10)`.

## Useful variations

- Floyd's tortoise and hare uses O(1) space: advance `slow` one step and `fast` two until they meet.
- Store `(state, step)` in a dict instead of a set to also learn the cycle's start and length.
- The same pattern finds repeated board states in Prison Cells After N Days and repeated remainders in Fraction to Recurring Decimal.

## Interview relevance

Happy Number is an easy problem whose real lesson is the stopping rule. Linked
List Cycle, repeated simulation states and recurring decimals all use a seen
set first, and interviewers often follow up with "can you do it in O(1) extra
space?" to lead you to Floyd's algorithm.
