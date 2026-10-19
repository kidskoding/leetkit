# Interval sweep with a heapq min-heap

Sort intervals by start, then keep a min-heap of the end times of the rooms
currently in use. The smallest end time sits at `ends[0]`, so checking
"is any room free yet?" is one comparison.

```python
import heapq

ends = []
for start, end in sorted(meetings):
    if ends and ends[0] <= start:
        heapq.heapreplace(ends, end)
    else:
        heapq.heappush(ends, end)
rooms = len(ends)
```

`sorted(meetings)` orders tuples by their first element, so no `key` is
needed. `heapreplace` pops the smallest end and pushes the new one in a single
O(log n) step, which is exactly "reuse the room that freed up earliest". The
clunky alternative compares each meeting against every open room, which is
O(n^2) and easy to get wrong at the boundaries.

## Useful variations

- `heapq.heappop(ends)` then `heapq.heappush(ends, end)` is the two-step form of `heapreplace`.
- Sort `starts` and `ends` separately and walk them with two pointers for a heap-free O(n log n) answer.
- Push `(end, room_id)` tuples when you must report which room each meeting gets.
- Change `<=` to `<` if a meeting ending at `t` must not share a room with one starting at `t`.

## Interview relevance

Meeting Rooms II, Minimum Number of Platforms, Car Pooling, and Single-Threaded
CPU all use this sweep. Interviewers watch for sorting by start, the boundary
rule at equal times, and knowing that `heapq` is a min-heap over a plain list.
