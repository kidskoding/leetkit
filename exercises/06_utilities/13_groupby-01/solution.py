from itertools import groupby


def compress_levels(levels):
    return [(level, len(list(run))) for level, run in groupby(levels)]
