from collections import defaultdict


def count_clusters(n, links):
    graph = defaultdict(list)
    for a, b in links:
        graph[a].append(b)
        graph[b].append(a)
    seen = set()
    clusters = 0
    for start in range(n):
        if start in seen:
            continue
        clusters += 1
        seen.add(start)
        stack = [start]
        while stack:
            node = stack.pop()
            for nxt in graph[node]:
                if nxt not in seen:
                    seen.add(nxt)
                    stack.append(nxt)
    return clusters
