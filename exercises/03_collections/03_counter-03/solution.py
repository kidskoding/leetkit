from collections import Counter


def can_construct(note, magazine):
    return not (Counter(note) - Counter(magazine))
