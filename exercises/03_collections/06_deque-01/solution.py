from collections import defaultdict, deque


def min_flights(n, flights, src, dst):
    graph = defaultdict(list)
    for a, b in flights:
        graph[a].append(b)
        graph[b].append(a)
    queue = deque([(src, 0)])
    visited = {src}
    while queue:
        city, count = queue.popleft()
        if city == dst:
            return count
        for nxt in graph[city]:
            if nxt not in visited:
                visited.add(nxt)
                queue.append((nxt, count + 1))
    return -1
