from starter import is_isomorphic


def test_isomorphic_words():
    assert is_isomorphic("paper", "title") is True
    assert is_isomorphic("egg", "add") is True


def test_one_character_maps_to_two():
    assert is_isomorphic("foo", "bar") is False


def test_two_characters_map_to_one():
    assert is_isomorphic("badc", "baba") is False
    assert is_isomorphic("ab", "aa") is False


def test_different_lengths():
    assert is_isomorphic("abc", "ab") is False


def test_empty_strings():
    assert is_isomorphic("", "") is True


def test_single_characters():
    assert is_isomorphic("x", "y") is True


def test_case_sensitive():
    assert is_isomorphic("aA", "bb") is False
    assert is_isomorphic("aA", "bB") is True
