import math


def max_profit(prices):
    min_price = math.inf
    best = 0
    for price in prices:
        min_price = min(min_price, price)
        best = max(best, price - min_price)
    return best
