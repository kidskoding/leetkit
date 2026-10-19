from bisect import insort


def running_medians(stream):
    window = []
    medians = []
    for value in stream:
        insort(window, value)
        mid = len(window) // 2
        if len(window) % 2:
            medians.append(float(window[mid]))
        else:
            medians.append((window[mid - 1] + window[mid]) / 2)
    return medians
