def find_redundant_cable(cables):
    parent = {}
    size = {}

    def find(x):
        parent.setdefault(x, x)
        size.setdefault(x, 1)
        while parent[x] != x:
            parent[x] = parent[parent[x]]
            x = parent[x]
        return x

    for cable in cables:
        a, b = cable
        ra, rb = find(a), find(b)
        if ra == rb:
            return cable
        if size[ra] < size[rb]:
            ra, rb = rb, ra
        parent[rb] = ra
        size[ra] += size[rb]
    return None
