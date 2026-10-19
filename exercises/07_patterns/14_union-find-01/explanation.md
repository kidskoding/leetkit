# Union-find with dicts

Union-find (disjoint set union) tracks which items belong to the same group.
`find(x)` walks parent links to the group's root; `union` points one root at
the other. Two dicts let it grow lazily for any hashable label, and
`setdefault` adds a new item as its own root the first time it is seen.

```python
parent, size = {}, {}

def find(x):
    parent.setdefault(x, x)
    size.setdefault(x, 1)
    while parent[x] != x:
        parent[x] = parent[parent[x]]
        x = parent[x]
    return x

ra, rb = find(a), find(b)
if ra != rb:
    if size[ra] < size[rb]:
        ra, rb = rb, ra
    parent[rb] = ra
    size[ra] += size[rb]
```

`parent[x] = parent[parent[x]]` is path halving, a one-line form of path
compression that keeps trees flat. Attaching the smaller tree under the larger
(union by size) bounds the height even before compression. Together each
operation is nearly O(1), and the iterative `while` loop never hits Python's
recursion limit. The alternative, rerunning BFS or DFS after every cable to
check connectivity, costs O(n) per cable.

## Useful variations

- `parent = list(range(n))` replaces the dicts when labels are exactly `0..n-1`.
- Recursive `parent[x] = find(parent[x])` is full path compression but recurses deeply on long chains.
- Count groups by starting at `n` and subtracting one on every successful union.
- Sort edges by weight first and union cheapest-first to get Kruskal's minimum spanning tree.

## Interview relevance

Redundant Connection, Number of Provinces, Graph Valid Tree, Accounts Merge,
and Kruskal's MST all use this structure. Interviewers check that you know
both optimisations, can explain near-constant amortised cost, and can adapt
the structure to string labels with dicts.
