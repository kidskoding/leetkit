from starter import group_shifted


def test_classic_example():
    words = ["abc", "bcd", "acef", "xyz", "az", "ba", "a", "z"]
    assert group_shifted(words) == [["abc", "bcd", "xyz"], ["acef"], ["az", "ba"], ["a", "z"]]


def test_empty_list():
    assert group_shifted([]) == []


def test_wraparound_matches():
    assert group_shifted(["za", "ab", "mn"]) == [["za", "ab", "mn"]]


def test_different_lengths_never_match():
    assert group_shifted(["a", "ab", "abc", "b"]) == [["a", "b"], ["ab"], ["abc"]]


def test_duplicates_are_kept():
    assert group_shifted(["ab", "ab", "cd"]) == [["ab", "ab", "cd"]]


def test_group_order_follows_first_appearance():
    assert group_shifted(["ace", "aa", "bdf", "zz", "q"]) == [["ace", "bdf"], ["aa", "zz"], ["q"]]
