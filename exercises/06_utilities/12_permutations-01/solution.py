from itertools import permutations


def arrangements(word):
    return sorted({"".join(p) for p in permutations(word)})
