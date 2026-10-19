from collections import Counter


def rank_words(words):
    counts = Counter(words)
    return sorted(counts, key=lambda w: (-counts[w], w))
