from starter import has_duplicate


def test_duplicate_present():
    assert has_duplicate([101, 205, 309, 205]) is True


def test_all_distinct():
    assert has_duplicate([1, 2, 3]) is False


def test_empty_list():
    assert has_duplicate([]) is False


def test_single_element():
    assert has_duplicate([42]) is False


def test_adjacent_duplicates():
    assert has_duplicate([7, 7]) is True


def test_negatives_and_zero():
    assert has_duplicate([-1, 0, 1, -1]) is True
    assert has_duplicate([-3, -2, -1, 0]) is False
