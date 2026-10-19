import math
from functools import reduce


def bundle_size(counts):
    return reduce(math.gcd, counts, 0)
