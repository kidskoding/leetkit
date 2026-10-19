from starter import count_in_range


def test_duplicates_on_both_bounds():
    assert count_in_range([5, 10, 10, 15, 20, 20, 30], 10, 20) == 5


def test_no_values_in_range():
    assert count_in_range([5, 10, 15], 11, 14) == 0


def test_range_covers_everything():
    assert count_in_range([2, 4, 6], 0, 100) == 3


def test_single_value_range():
    assert count_in_range([1, 3, 3, 3, 7], 3, 3) == 3


def test_range_outside_list():
    assert count_in_range([10, 20], 30, 40) == 0
    assert count_in_range([10, 20], -5, 5) == 0


def test_lo_greater_than_hi():
    assert count_in_range([1, 2, 3, 4], 3, 2) == 0


def test_empty_list():
    assert count_in_range([], 1, 10) == 0


def test_negative_prices():
    assert count_in_range([-8, -3, -3, 0, 4], -3, 0) == 3
