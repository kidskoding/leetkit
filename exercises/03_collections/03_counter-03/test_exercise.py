from starter import can_construct


def test_enough_letters():
    assert can_construct("aa", "aab") is True


def test_not_enough_copies():
    assert can_construct("aa", "ab") is False


def test_missing_letter():
    assert can_construct("abc", "aabbdd") is False


def test_exact_match_in_any_order():
    assert can_construct("listen", "silent") is True


def test_empty_note():
    assert can_construct("", "anything") is True


def test_empty_magazine():
    assert can_construct("a", "") is False


def test_magazine_with_many_extras():
    assert can_construct("cab", "zzzbbbaaaccc") is True
