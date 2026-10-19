# Generator expressions

A generator expression looks like a list comprehension with round brackets
instead of square ones. It produces values lazily, one at a time.

```python
total = sum(price * qty for price, qty in items)
```

When the generator is the only argument to a function, the extra parentheses
can be dropped, as above. Compare with the list version:

```python
total = sum([price * qty for price, qty in items])
```

Both give the same answer, but the list version builds the whole list in
memory first and then throws it away. The generator feeds `sum` directly.

## Useful variations

- `max(len(w) for w in words)` finds the longest word length.
- `max((len(w) for w in words), default=0)` handles an empty input; the
  generator needs its own parentheses once there is a second argument.
- `sum(1 for x in nums if x > 0)` counts matches without building a list.
- `next((x for x in nums if x % 7 == 0), None)` returns the first match or a
  default.

## Interview relevance

Totals, counts, minimums and maximums appear in almost every problem: sum of a
window, number of valid pairs, best score so far. Feeding a generator into
`sum`, `max`, `min`, `any` or `all` keeps that code to one line and uses O(1)
extra memory, which is worth mentioning when you discuss complexity.
