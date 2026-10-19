from starter import find_kth_largest


def test_second_largest():
    assert find_kth_largest([3, 2, 1, 5, 6, 4], 2) == 5


def test_duplicates_count_separately():
    assert find_kth_largest([3, 2, 3, 1, 2, 4, 5, 5, 6], 4) == 4
    assert find_kth_largest([5, 5, 3], 2) == 5


def test_k_equals_one_is_maximum():
    assert find_kth_largest([7, 10, 4, 3, 20, 15], 1) == 20


def test_k_equals_length_is_minimum():
    assert find_kth_largest([7, 10, 4, 3, 20, 15], 6) == 3


def test_single_element():
    assert find_kth_largest([42], 1) == 42


def test_negative_numbers():
    assert find_kth_largest([-1, -5, -3, -2], 3) == -3
