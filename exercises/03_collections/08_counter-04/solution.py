from collections import Counter


def find_anagrams(s, p):
    k = len(p)
    target = Counter(p)
    window = Counter()
    starts = []
    for i, ch in enumerate(s):
        window[ch] += 1
        if i >= k:
            window[s[i - k]] -= 1
        if window == target:
            starts.append(i - k + 1)
    return starts
