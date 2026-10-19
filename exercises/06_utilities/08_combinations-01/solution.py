from itertools import combinations


def pairs_with_sum(prices, target):
    return [(a, b) for a, b in combinations(prices, 2) if a + b == target]
