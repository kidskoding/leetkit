# DFS with an adjacency list and a stack

Graph problems usually hand you an edge list. Convert it to an adjacency list
first, then run a depth-first search from every node you have not visited yet;
each such start is a new component.

```python
from collections import defaultdict

graph = defaultdict(list)
for a, b in links:
    graph[a].append(b)
    graph[b].append(a)

stack = [start]
while stack:
    node = stack.pop()
    for nxt in graph[node]:
        if nxt not in seen:
            seen.add(nxt)
            stack.append(nxt)
```

`defaultdict(list)` removes the "create the list if missing" branch, and
`graph[node]` on an isolated node just returns an empty list. Tuple unpacking in
`for a, b in links` reads each edge directly. A plain list with `append` and
`pop` is a stack, so the search needs no recursion.

## Useful variations

- Recursive `def dfs(node):` is shorter, but Python stops at about 1000 frames by default.
- `graph = [[] for _ in range(n)]` works when nodes are exactly `0..n-1`.
- Swap the stack for `deque` and `popleft()` to get BFS with the same code.
- Union-find with a `parent` list answers the same question without an adjacency list.

## Interview relevance

Connected components, number of provinces, graph valid tree, course schedule,
and clone graph all begin with "edge list to adjacency list". Interviewers check
that you handle isolated nodes, avoid revisiting, and know when recursion depth
becomes a risk.
