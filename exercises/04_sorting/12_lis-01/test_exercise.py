from starter import length_of_lis


def test_classic_example():
    assert length_of_lis([10, 9, 2, 5, 3, 7, 101, 18]) == 4


def test_repeated_values_in_sequence():
    assert length_of_lis([0, 1, 0, 3, 2, 3]) == 4


def test_equal_values_are_not_increasing():
    assert length_of_lis([7, 7, 7, 7]) == 1
    assert length_of_lis([1, 2, 2, 3]) == 3


def test_strictly_decreasing():
    assert length_of_lis([5, 4, 3, 2, 1]) == 1


def test_negatives():
    assert length_of_lis([-3, -1, -2, 0, 5, -10]) == 4


def test_empty_and_single():
    assert length_of_lis([]) == 0
    assert length_of_lis([42]) == 1


def test_long_increasing_run():
    assert length_of_lis(list(range(1000))) == 1000
