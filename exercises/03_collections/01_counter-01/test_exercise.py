from starter import is_anagram


def test_simple_anagram():
    assert is_anagram("listen", "silent") is True


def test_not_an_anagram():
    assert is_anagram("rat", "car") is False


def test_same_letters_different_counts():
    assert is_anagram("aab", "abb") is False


def test_different_lengths():
    assert is_anagram("abc", "abcc") is False


def test_case_sensitive():
    assert is_anagram("Dog", "god") is False


def test_empty_strings():
    assert is_anagram("", "") is True


def test_repeated_letters():
    assert is_anagram("anagram", "nagaram") is True
