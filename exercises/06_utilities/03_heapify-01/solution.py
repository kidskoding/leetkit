import heapq


def k_smallest(prices, k):
    heap = list(prices)
    heapq.heapify(heap)
    return [heapq.heappop(heap) for _ in range(min(k, len(heap)))]
