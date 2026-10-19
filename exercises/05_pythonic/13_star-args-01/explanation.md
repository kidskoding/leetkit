# *args and **kwargs

In a function definition, `*name` collects any extra positional arguments into
a tuple and `**name` collects any extra keyword arguments into a dict:

```python
def merge_settings(*layers, **overrides):
    merged = {}
    for layer in layers:
        merged.update(layer)
    merged.update(overrides)
    return merged
```

At a call site the stars do the opposite: `*` spreads a sequence into
positional arguments and `**` spreads a dict into keyword arguments.

```python
configs = [defaults, from_file]
merge_settings(*configs, **cli_flags)
# same as merge_settings(defaults, from_file, verbose=True, ...)
```

The alternative, `def merge_settings(layers, overrides)`, forces every caller
to build a list and a dict first, even for one layer.

Starting from a fresh `{}` matters: calling `layers[0].update(...)` would
quietly change the caller's dict.

## Useful variations

- `{**a, **b}` merges two dicts into a new one, with `b` winning; `a | b` does the same in 3.9+.
- `def wrapper(*args, **kwargs): return func(*args, **kwargs)` forwards everything in a decorator.
- `print(*row)` prints a list's items separated by spaces.
- A bare `*` in a signature, as in `def f(a, *, key)`, makes the following parameters keyword-only.

## Interview relevance

Decorators, memoization wrappers and small utility helpers all need variadic
signatures, and "what do `*args` and `**kwargs` do?" is a standard Python
screening question. Being able to use both sides, defining and unpacking,
shows you understand how Python binds arguments rather than just
recognizing the syntax.
