from starter import count_jewels


def test_mixed_case_example():
    assert count_jewels("aA", "aAAbbbb") == 3


def test_case_sensitive():
    assert count_jewels("z", "ZZ") == 0


def test_all_stones_are_jewels():
    assert count_jewels("abc", "cabbac") == 6


def test_empty_stones():
    assert count_jewels("abc", "") == 0


def test_empty_jewels():
    assert count_jewels("", "abc") == 0


def test_digits_and_symbols():
    assert count_jewels("1#", "1a#b#11") == 5
