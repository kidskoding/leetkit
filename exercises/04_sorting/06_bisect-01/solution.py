from bisect import bisect_left


def search_insert(nums, target):
    return bisect_left(nums, target)
