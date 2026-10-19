# Neighbour pairs with zip(xs, xs[1:])

Once a list is sorted, many questions only need each element compared with the
next one. `zip(xs, xs[1:])` yields those `(previous, current)` pairs, and
`all(...)` turns the per-pair rule into one boolean.

```python
meetings = sorted(intervals, key=lambda iv: iv[0])
ok = all(prev[1] <= cur[0] for prev, cur in zip(meetings, meetings[1:]))
```

The index-based version needs `range(1, len(meetings))` and
`meetings[i - 1][1] <= meetings[i][0]`, which is easy to get off by one.
Sorting first is what makes the neighbour check enough: if meeting A starts
before B and they overlap, then every meeting between them in sorted order
overlaps A too, so some adjacent pair always clashes.

## Useful variations

- `itertools.pairwise(meetings)` (Python 3.10+) yields the same pairs without slicing.
- `any(prev[1] > cur[0] for ...)` asks the opposite question: is there a conflict?
- `sorted(intervals)` also works here, since lists compare by start first.
- `all(...)` on an empty generator is `True`, which handles zero or one meeting for free.

## Interview relevance

Meeting Rooms, checking that a list is sorted, and finding the smallest gap
between sorted values all reduce to "sort, then look at neighbours". The
`zip(xs, xs[1:])` idiom keeps the pass to one readable line and leaves time
for the follow-up, Meeting Rooms II, which asks how many rooms you need.
