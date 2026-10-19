def pair_with_sum(prices, target):
    lo, hi = 0, len(prices) - 1
    while lo < hi:
        total = prices[lo] + prices[hi]
        if total == target:
            return (lo, hi)
        if total < target:
            lo += 1
        else:
            hi -= 1
    return (-1, -1)
