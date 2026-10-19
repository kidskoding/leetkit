import math


def min_coins(coins, amount):
    best = [0] + [math.inf] * amount
    for a in range(1, amount + 1):
        for coin in coins:
            if coin <= a:
                best[a] = min(best[a], best[a - coin] + 1)
    return -1 if best[amount] == math.inf else best[amount]
