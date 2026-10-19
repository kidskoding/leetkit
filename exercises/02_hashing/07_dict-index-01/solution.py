def first_unique_index(s):
    positions = {}
    for i, ch in enumerate(s):
        positions[ch] = -1 if ch in positions else i
    return next((i for i in positions.values() if i != -1), -1)
