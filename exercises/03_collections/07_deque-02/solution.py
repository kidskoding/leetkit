from collections import deque


def moving_average(readings, k):
    window = deque(maxlen=k)
    averages = []
    for reading in readings:
        window.append(reading)
        averages.append(sum(window) / len(window))
    return averages
