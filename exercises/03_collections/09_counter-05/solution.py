from collections import Counter


def common_chars(words):
    common = Counter(words[0])
    for word in words[1:]:
        common &= Counter(word)
    return sorted(common.elements())
