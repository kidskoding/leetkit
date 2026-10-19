from starter import count_periods


def test_overlapping_runs():
    assert count_periods([1, 1, 1], 2) == 2


def test_single_element_and_longer_run():
    assert count_periods([1, 2, 3], 3) == 2


def test_negative_values():
    assert count_periods([1, -1, 0], 0) == 3
    assert count_periods([2, -2, 2], 2) == 3


def test_all_zeros_target_zero():
    assert count_periods([0, 0, 0], 0) == 6


def test_no_matching_run():
    assert count_periods([1, 2], 10) == 0


def test_empty_list():
    assert count_periods([], 0) == 0


def test_whole_list_only():
    assert count_periods([5, 2, -3], 4) == 1
