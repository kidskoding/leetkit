import heapq


def top_orders(orders, k):
    return [order["id"] for order in heapq.nlargest(k, orders, key=lambda o: o["total"])]
