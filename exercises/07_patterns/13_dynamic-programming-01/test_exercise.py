from starter import min_coins


def test_standard_coins():
    assert min_coins([1, 2, 5], 11) == 3


def test_greedy_is_not_optimal():
    assert min_coins([1, 3, 4], 6) == 2


def test_impossible_amount():
    assert min_coins([2], 3) == -1
    assert min_coins([5, 10], 7) == -1


def test_zero_amount_needs_no_coins():
    assert min_coins([1, 2], 0) == 0
    assert min_coins([], 0) == 0


def test_no_coins_positive_amount():
    assert min_coins([], 5) == -1


def test_unsorted_coins():
    assert min_coins([10, 1, 5, 2], 27) == 4


def test_larger_amount():
    assert min_coins([3, 7, 11], 100) == 12
    assert min_coins([1], 250) == 250
