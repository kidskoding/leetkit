# defaultdict for adjacency lists

An adjacency list maps each node to its neighbours. `defaultdict` creates the
empty container the first time a node is touched, so building one is a
two-line loop.

```python
from collections import defaultdict

graph = defaultdict(set)
for a, b in edges:
    graph[a].add(b)
    graph[b].add(a)
```

Without it, every edge needs `if a not in graph: graph[a] = set()` twice.
Using `set` as the factory also removes duplicate edges automatically.

## Useful variations

- `defaultdict(list)` keeps neighbours in edge order and is the usual choice when edges are unique.
- For a directed graph, add only `graph[a].append(b)`.
- For weighted edges, `graph[a].append((b, weight))` stores the cost alongside the neighbour.
- Reading `graph[x]` for an unknown node inserts it; use `graph.get(x, ())` when you only want to look.

## Interview relevance

Graph questions rarely hand you an adjacency list; they give edges, a list of
prerequisites, or word pairs. Building the graph quickly and correctly
(remembering both directions for undirected edges) is the first step of BFS,
DFS, topological sort, and union-find discussions.
