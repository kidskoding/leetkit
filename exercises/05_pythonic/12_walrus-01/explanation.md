# The walrus operator :=

An assignment expression `name := value` stores a value and also evaluates to
it, so you can bind a result in the middle of a condition and use it again:

```python
return [
    (word, score)
    for word in words
    if (score := word_score(word)) > cutoff
]
```

Without it you compute the transform twice:

```python
[(w, word_score(w)) for w in words if word_score(w) > cutoff]
```

or give up on the comprehension and write a loop with a temporary variable.

The parentheses matter: `score := f(w) > cutoff` would bind the boolean
result of the comparison, not the score. Names bound with `:=` inside a
comprehension leak into the enclosing function scope, which is occasionally
useful and occasionally surprising.

## Useful variations

- `while (line := stream.readline()):` reads until an empty string.
- `if (m := re.match(pattern, s)):` tests and keeps a match object in one step.
- `while (n := n // 10):` repeats until a running value hits zero.
- `[y for x in xs if (y := f(x)) is not None]` keeps only successful transforms.

## Interview relevance

"Map then filter on the mapped value" is everywhere: scores above a threshold,
parsed values that are valid, distances inside a radius. Using `:=` keeps that
in a single comprehension without doing the work twice, and it is a modern
Python feature (3.8+) that interviewers notice when used with restraint.
