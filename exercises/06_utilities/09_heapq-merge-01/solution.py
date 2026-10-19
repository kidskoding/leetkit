import heapq


def merge_streams(streams):
    return list(heapq.merge(*streams))
