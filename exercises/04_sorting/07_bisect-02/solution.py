from bisect import bisect_left, bisect_right


def count_in_range(prices, lo, hi):
    return max(0, bisect_right(prices, hi) - bisect_left(prices, lo))
