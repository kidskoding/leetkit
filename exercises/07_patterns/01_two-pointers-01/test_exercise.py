from starter import pair_with_sum


def test_pair_at_start():
    assert pair_with_sum([2, 7, 11, 15], 9) == (0, 1)


def test_pair_in_middle():
    assert pair_with_sum([1, 3, 4, 6, 12], 10) == (2, 3)


def test_pair_at_both_ends():
    assert pair_with_sum([1, 4, 5, 8, 20], 21) == (0, 4)


def test_duplicate_prices():
    assert pair_with_sum([3, 3, 8], 6) == (0, 1)


def test_negative_values():
    assert pair_with_sum([-4, -1, 2, 7], 1) == (1, 2)


def test_no_pair():
    assert pair_with_sum([5, 10], 4) == (-1, -1)


def test_empty_and_single():
    assert pair_with_sum([], 0) == (-1, -1)
    assert pair_with_sum([5], 10) == (-1, -1)
