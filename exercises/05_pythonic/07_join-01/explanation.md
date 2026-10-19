# str.join()

`sep.join(iterable)` concatenates an iterable of strings, putting `sep`
between each pair.

```python
"".join(["a3", "b1", "c2"])      # "a3b1c2"
", ".join(["red", "green"])      # "red, green"
```

The idiom for building a string piece by piece is to collect the pieces in a
list and join once at the end:

```python
parts = []
for ch, count in runs:
    parts.append(f"{ch}{count}")
return "".join(parts)
```

Compare with `result += f"{ch}{count}"` inside the loop. Strings are
immutable, so each `+=` may copy everything built so far, which can make the
loop quadratic. `join` computes the final size once and copies each piece once.

`join` only accepts strings: `",".join([1, 2])` raises TypeError. Convert first
with `",".join(map(str, nums))`.

## Useful variations

- `" ".join(reversed(s.split()))` reverses the words in a sentence.
- `"\n".join(lines)` builds multi-line output.
- `"".join(ch for ch in s if ch.isalnum())` filters characters with a generator.
- `"".join(f"{k}{len(list(g))}" for k, g in itertools.groupby(s))` is a compact
  run-length encoder once you know `groupby`.

## Interview relevance

String Compression, Count and Say, Reverse Words in a String and tree
serialization all end by assembling a string. Using a list plus `join` keeps
the solution linear, and saying why (immutable strings, repeated copying) is an
easy signal of Python fluency.
