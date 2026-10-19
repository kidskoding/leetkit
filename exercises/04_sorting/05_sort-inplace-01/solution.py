def top_k(scores, k):
    return sorted(scores, reverse=True)[:k]
