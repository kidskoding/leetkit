from starter import longest_consecutive


def test_unsorted_run():
    assert longest_consecutive([100, 4, 200, 1, 3, 2]) == 4


def test_empty_list():
    assert longest_consecutive([]) == 0


def test_single_element():
    assert longest_consecutive([7]) == 1


def test_duplicates_count_once():
    assert longest_consecutive([1, 2, 2, 3]) == 3
    assert longest_consecutive([0, 3, 7, 2, 5, 8, 4, 6, 0, 1]) == 9


def test_negative_numbers():
    assert longest_consecutive([-1, -3, 5, -2, 0]) == 4


def test_no_neighbours():
    assert longest_consecutive([10, 30, 20]) == 1


def test_picks_longer_of_two_runs():
    assert longest_consecutive([50, 1, 51, 2, 52, 53, 3]) == 4
