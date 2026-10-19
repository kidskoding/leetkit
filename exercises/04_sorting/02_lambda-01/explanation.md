# lambda as a sort key

A `lambda` is a one-expression function written inline. It is the usual way
to pull one field out of a record for sorting.

```python
intervals = [[8, 10], [1, 3], [2, 6]]
sorted(intervals, key=lambda iv: iv[0])
# [[1, 3], [2, 6], [8, 10]]
```

Without a key, lists compare element by element, which happens to sort by
start here too, but the lambda states the intent and works for records where
the field you want is not first.

Once sorted, Merge Intervals is a single sweep: each interval either overlaps
the last merged one (extend its end) or starts a new group.

```python
for start, end in sorted(intervals, key=lambda iv: iv[0]):
    if merged and start <= merged[-1][1]:
        merged[-1][1] = max(merged[-1][1], end)
    else:
        merged.append([start, end])
```

## Useful variations

- `key=lambda iv: iv[1]` sorts by end time, the right order for greedy
  "maximum non-overlapping intervals" problems.
- `from operator import itemgetter` then `key=itemgetter(0)` does the same as
  the lambda and is slightly faster.
- `key=lambda p: p.age` works on objects as well as lists and tuples.

## Interview relevance

Interval problems (Merge Intervals, Meeting Rooms, Insert Interval, Minimum
Arrows) almost all begin with sorting by start or end. Writing the lambda key
without hesitation, and appending a fresh `[start, end]` list so you never
mutate the caller's data, signals fluency.
