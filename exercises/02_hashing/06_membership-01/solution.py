def count_jewels(jewels, stones):
    jewel_set = set(jewels)
    return sum(1 for ch in stones if ch in jewel_set)
