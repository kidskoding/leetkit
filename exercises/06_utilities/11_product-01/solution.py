from itertools import product


def all_codes(alphabet, n):
    return ["".join(chars) for chars in product(alphabet, repeat=n)]
