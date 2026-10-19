from collections import Counter


def count_periods(changes, k):
    seen = Counter({0: 1})
    total = count = 0
    for x in changes:
        total += x
        count += seen[total - k]
        seen[total] += 1
    return count
