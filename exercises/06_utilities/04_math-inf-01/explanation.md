# math.inf as a sentinel

`math.inf` is a float larger than every other number, and `-math.inf` is
smaller than every other number.

```python
import math

min_price = math.inf
for price in prices:
    min_price = min(min_price, price)
```

The first comparison always replaces the sentinel, so there is no special case
for the first element. Compare the clunky alternatives:

```python
min_price = prices[0]      # crashes on an empty list
min_price = 10**9          # wrong if a price is ever larger
```

## Useful variations

- `best = -math.inf` starts a running maximum, e.g. Maximum Subarray.
- `dist = [math.inf] * n` initialises unreached nodes in Dijkstra or Bellman-Ford.
- `float("inf")` is the same value without the import.
- Check `if answer == math.inf: return -1` to report "no solution" in DP problems like Coin Change.

## Interview relevance

Sentinels show up in running min/max loops, shortest-path distance arrays, and
minimum-cost DP tables. `math.inf` reads clearly, never overflows, and compares
correctly with ints, so it is the idiomatic Python choice.
