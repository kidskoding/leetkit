def first_bad_build(n, is_bad):
    lo, hi = 1, n + 1
    while lo < hi:
        mid = (lo + hi) // 2
        if is_bad(mid):
            hi = mid
        else:
            lo = mid + 1
    return lo if lo <= n else -1
