from starter import max_sliding_window


def test_leetcode_example():
    assert max_sliding_window([1, 3, -1, -3, 5, 3, 6, 7], 3) == [3, 3, 5, 5, 6, 7]


def test_window_of_one_returns_input():
    assert max_sliding_window([4, -2, 7, 0], 1) == [4, -2, 7, 0]


def test_window_covers_whole_list():
    assert max_sliding_window([2, 9, 4, 1], 4) == [9]


def test_decreasing_values_expire_from_front():
    assert max_sliding_window([9, 8, 7, 6], 2) == [9, 8, 7]


def test_increasing_values():
    assert max_sliding_window([1, 2, 3, 4, 5], 3) == [3, 4, 5]


def test_duplicates_and_negatives():
    assert max_sliding_window([-1, -1, -3, -1, -5], 2) == [-1, -1, -1, -1]


def test_empty_input():
    assert max_sliding_window([], 3) == []
