# Topological sort with Kahn's algorithm

Kahn's algorithm repeatedly takes a node with no remaining incoming edges.
`defaultdict(list)` builds the adjacency list without "create if missing"
checks, a dict counts each node's in-degree, and a queue holds the nodes that
are ready.

```python
unlocks = defaultdict(list)
indegree = {c: 0 for c in range(n)}
for course, before in prereqs:
    unlocks[before].append(course)
    indegree[course] += 1

ready = [c for c in range(n) if indegree[c] == 0]
heapq.heapify(ready)
while ready:
    course = heapq.heappop(ready)
    order.append(course)
    for nxt in unlocks[course]:
        indegree[nxt] -= 1
        if indegree[nxt] == 0:
            heapq.heappush(ready, nxt)
```

Courses on a cycle never reach in-degree 0, so they never enter the heap.
Checking `len(order) == n` at the end is the whole cycle test; no separate
DFS with colour marks is needed. The heap makes "smallest ready course first"
automatic and costs only an extra log factor.

## Useful variations

- When any valid order is acceptable, use `deque(ready)` with `popleft()` and `append()` instead of the heap.
- `graphlib.TopologicalSorter(graph).static_order()` (Python 3.9+) does the whole job and raises `CycleError` on cycles.
- Return `len(order) == n` alone for Course Schedule I, which only asks whether finishing is possible.
- For Alien Dictionary, build the edges by comparing adjacent words with `zip(w1, w2)`.

## Interview relevance

Course Schedule, Alien Dictionary, Parallel Courses, and build-dependency
ordering all use this loop. Interviewers look for correct in-degree
bookkeeping, a clean cycle check, and a deterministic order when they ask for
one; knowing that a heap gives the lexicographically smallest order is a
common follow-up.
