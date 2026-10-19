from starter import two_sum


def test_pair_at_start():
    assert two_sum([2, 7, 11, 15], 9) == [0, 1]


def test_pair_not_adjacent():
    assert two_sum([1, 5, 8, 3], 4) == [0, 3]


def test_equal_values():
    assert two_sum([3, 3], 6) == [0, 1]


def test_does_not_reuse_same_index():
    assert two_sum([3, 2, 4], 6) == [1, 2]


def test_negative_numbers():
    assert two_sum([-4, 10, -1, 6], 2) == [0, 3]


def test_zero_target():
    assert two_sum([5, 0, 9, 0], 0) == [1, 3]
