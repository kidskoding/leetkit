# map() and filter()

`map(func, iterable)` applies `func` to every item, and `filter(pred, iterable)`
keeps the items for which `pred` returns a truthy value. Both return lazy
iterators, so wrap them in `list()` when you need a list:

```python
return list(filter(lambda n: n > 0, map(int, tokens)))
```

They shine when the function already exists. `map(int, tokens)` is shorter
than `[int(t) for t in tokens]` and reads as "convert everything".

Once you need a `lambda`, the comprehension is usually clearer:

```python
nums = map(int, tokens)
return [n for n in nums if n > 0]
```

Both are fine in an interview; pick the one that reads best and be ready to
explain the other.

## Useful variations

- `list(map(int, input().split()))` is the standard way to read a line of numbers.
- `filter(None, items)` drops falsy values such as `""`, `0` and `None`.
- `map(str, nums)` feeds `",".join(...)`, which only accepts strings.
- `map(max, xs, ys)` walks several iterables in parallel, like `zip`.

## Interview relevance

Reading and cleaning input is the first step of many problems, and
`map(int, ...)` is the idiom everyone expects. Questions like "what does map
return in Python 3?" (a lazy iterator, not a list) or "map or comprehension?"
come up in Python screens; knowing both styles and their trade-offs is an easy
way to show fluency.
