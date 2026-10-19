def group_by_letters(words):
    groups = {}
    for word in words:
        groups.setdefault(frozenset(word), []).append(word)
    return list(groups.values())
