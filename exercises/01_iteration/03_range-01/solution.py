def sample_from_end(readings, k):
    return [readings[i] for i in range(len(readings) - 1, -1, -k)]
