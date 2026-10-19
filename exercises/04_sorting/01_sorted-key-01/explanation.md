# sorted() with key=

`sorted(iterable, key=func)` returns a new list ordered by `func(item)`
instead of by the items themselves.

```python
sorted(["banana", "fig", "apple"], key=len)
# ['fig', 'apple', 'banana']
```

The key is called once per element, and the results are compared. Any
one-argument function works, including built-ins like `len`, `abs` and
`str.lower`.

Python's sort is stable: elements whose keys are equal stay in the order they
had in the input. That is why `["bb", "aa"]` stays `["bb", "aa"]` when sorted
by length, with no tie-breaking code at all.

## Useful variations

- `sorted(names, key=str.lower)` sorts case-insensitively.
- `sorted(nums, key=abs)` orders by distance from zero.
- `sorted(words, key=len, reverse=True)` puts the longest first; ties still
  keep their input order.
- `words.sort(key=len)` does the same thing in place and returns `None`.

## Interview relevance

Many problems start with "sort by X": strings by length, intervals by start,
tasks by deadline. A key function expresses that in one call and avoids
writing a comparator. When an interviewer asks how ties are handled, the answer
"Python's sort is stable, so they keep input order" is precise and correct.
