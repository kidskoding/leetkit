# Conditional expressions

A conditional expression chooses between two values in a single expression:

```python
label = "even" if n % 2 == 0 else "odd"
```

Because it is an expression, it can sit in the output position of a
comprehension, and it can be chained for more than two cases:

```python
[
    "FizzBuzz" if i % 15 == 0 else "Fizz" if i % 3 == 0 else "Buzz" if i % 5 == 0 else str(i)
    for i in range(1, n + 1)
]
```

The conditions are tested left to right, so put the most specific one (15)
first.

Do not confuse the two places `if` can go in a comprehension. An `if` after
the `for` filters items out; an `if ... else` before the `for` chooses a value
for every item:

```python
[x for x in nums if x > 0]          # drops non-positive values
[x if x > 0 else 0 for x in nums]   # keeps every position, replaces negatives
```

## Useful variations

- `[min(max(x, lo), hi) for x in nums]` is often clearer than a chain for clamping.
- `sign = 1 if x >= 0 else -1` keeps arithmetic code compact.
- `best = a if a > b else b` reads well, but `max(a, b)` says it better.
- Beyond two or three branches, a helper function or a dict lookup is clearer.

## Interview relevance

Conditional expressions keep per-element decisions inline: labelling cells in
a grid, mapping characters, choosing a direction. FizzBuzz itself is still a
common screening question, and writing it as one comprehension shows you know
the expression form. Know when to stop chaining, though; readability counts
when you are explaining code out loud.
