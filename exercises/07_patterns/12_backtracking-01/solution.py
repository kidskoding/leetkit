def crate_loads(weights, target):
    weights = sorted(weights)
    loads = []
    path = []

    def backtrack(start, remaining):
        if remaining == 0:
            loads.append(path[:])
            return
        for i in range(start, len(weights)):
            if weights[i] > remaining:
                break
            path.append(weights[i])
            backtrack(i, remaining - weights[i])
            path.pop()

    backtrack(0, target)
    return loads
