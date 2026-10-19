def high_scoring_words(words, cutoff):
    return [
        (word, score)
        for word in words
        if (score := sum(ord(ch) - ord("a") + 1 for ch in word)) > cutoff
    ]
