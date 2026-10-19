from starter import group_anagrams


def test_classic_example():
    words = ["eat", "tea", "tan", "ate", "nat", "bat"]
    assert group_anagrams(words) == [["eat", "tea", "ate"], ["tan", "nat"], ["bat"]]


def test_groups_in_first_seen_order():
    assert group_anagrams(["dog", "cat", "god", "act"]) == [["dog", "god"], ["cat", "act"]]


def test_no_anagrams():
    assert group_anagrams(["abc", "xyz"]) == [["abc"], ["xyz"]]


def test_duplicates_are_kept():
    assert group_anagrams(["ab", "ba", "ab"]) == [["ab", "ba", "ab"]]


def test_empty_string_word():
    assert group_anagrams(["", "a", ""]) == [["", ""], ["a"]]


def test_same_letters_different_counts():
    assert group_anagrams(["aab", "abb"]) == [["aab"], ["abb"]]


def test_empty_input():
    assert group_anagrams([]) == []
