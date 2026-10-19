from collections import defaultdict


def group_shifted(words):
    groups = defaultdict(list)
    for w in words:
        key = tuple((ord(b) - ord(a)) % 26 for a, b in zip(w, w[1:]))
        groups[key].append(w)
    return list(groups.values())
