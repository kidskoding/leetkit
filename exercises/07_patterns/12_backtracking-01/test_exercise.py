from starter import crate_loads


def test_two_loads():
    assert crate_loads([2, 3, 6, 7], 7) == [[2, 2, 3], [7]]


def test_unsorted_weights_lexicographic_output():
    assert crate_loads([5, 3, 2], 8) == [[2, 2, 2, 2], [2, 3, 3], [3, 5]]


def test_no_load_fits():
    assert crate_loads([4, 6], 5) == []


def test_empty_weights():
    assert crate_loads([], 3) == []


def test_single_weight_repeated():
    assert crate_loads([1], 3) == [[1, 1, 1]]


def test_single_item_hits_target_too():
    assert crate_loads([10, 5], 10) == [[5, 5], [10]]


def test_many_loads():
    assert crate_loads([1, 2, 3], 4) == [[1, 1, 1, 1], [1, 1, 2], [1, 3], [2, 2]]
