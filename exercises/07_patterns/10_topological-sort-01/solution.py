import heapq
from collections import defaultdict


def course_order(n, prereqs):
    unlocks = defaultdict(list)
    indegree = {course: 0 for course in range(n)}
    for course, before in prereqs:
        unlocks[before].append(course)
        indegree[course] += 1
    ready = [course for course in range(n) if indegree[course] == 0]
    heapq.heapify(ready)
    order = []
    while ready:
        course = heapq.heappop(ready)
        order.append(course)
        for nxt in unlocks[course]:
            indegree[nxt] -= 1
            if indegree[nxt] == 0:
                heapq.heappush(ready, nxt)
    return order if len(order) == n else []
