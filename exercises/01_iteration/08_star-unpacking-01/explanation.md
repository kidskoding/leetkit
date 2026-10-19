# Extended (starred) unpacking

A single starred name in an assignment target soaks up "everything else" as a
list, while the plain names take items from the ends.

```python
origin, *stops, destination = "LAX-ORD-DEN-JFK".split("-")
# origin = "LAX", stops = ["ORD", "DEN"], destination = "JFK"

first, *rest = [3, 1, 4]      # first = 3, rest = [1, 4]
*init, last = [3, 1, 4]       # init = [3, 1], last = 4
```

The starred name is always a list, even when it captures nothing:

```python
a, *middle, b = ["x", "y"]    # middle = []
```

Compare with manual slicing:

```python
parts = route.split("-")
origin, stops, destination = parts[0], parts[1:-1], parts[-1]
```

## Useful variations

- `cmd, *args = line.split()` separates a command from its arguments.
- `head, *_ = rows` takes the first item and ignores the rest.
- `for name, *scores in records:` unpacks inside a loop.
- Only one starred name is allowed per target, and the plain names must be
  satisfiable: `a, *b, c = [1]` raises `ValueError`.

## Interview relevance

Parsing problems (log lines, commands, CSV-like rows) and recursion on
"first and rest" become one clear line with starred unpacking. It also avoids
the off-by-one slips that creep into `parts[1:-1]`-style slicing.
