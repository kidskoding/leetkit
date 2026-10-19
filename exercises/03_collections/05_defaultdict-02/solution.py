from collections import defaultdict


def build_graph(edges):
    graph = defaultdict(set)
    for a, b in edges:
        graph[a].add(b)
        graph[b].add(a)
    return {user: sorted(friends) for user, friends in graph.items()}
