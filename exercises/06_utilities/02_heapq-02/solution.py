import heapq


def process_tasks(tasks):
    heap = []
    for i, (name, priority) in enumerate(tasks):
        heapq.heappush(heap, (priority, i, name))
    order = []
    while heap:
        _, _, name = heapq.heappop(heap)
        order.append(name)
    return order
