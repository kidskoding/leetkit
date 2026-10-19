from starter import find_anagrams


def test_two_separate_matches():
    assert find_anagrams("cbaebabacd", "abc") == [0, 6]


def test_overlapping_matches():
    assert find_anagrams("abab", "ab") == [0, 1, 2]


def test_repeated_letters_need_matching_counts():
    assert find_anagrams("aaabaa", "aab") == [1, 2, 3]


def test_no_match():
    assert find_anagrams("hello", "xyz") == []


def test_pattern_longer_than_string():
    assert find_anagrams("ab", "abc") == []


def test_whole_string_is_an_anagram():
    assert find_anagrams("listen", "silent") == [0]


def test_single_letter_pattern():
    assert find_anagrams("banana", "a") == [1, 3, 5]
