from collections import deque


def max_sliding_window(nums, k):
    window = deque()
    maxima = []
    for i, x in enumerate(nums):
        while window and nums[window[-1]] <= x:
            window.pop()
        window.append(i)
        if window[0] <= i - k:
            window.popleft()
        if i >= k - 1:
            maxima.append(nums[window[0]])
    return maxima
