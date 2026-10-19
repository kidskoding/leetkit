from itertools import accumulate


def range_sums(nums, queries):
    prefix = list(accumulate(nums, initial=0))
    return [prefix[right + 1] - prefix[left] for left, right in queries]
