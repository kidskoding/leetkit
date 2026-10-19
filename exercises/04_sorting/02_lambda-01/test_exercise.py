from starter import merge_intervals


def test_overlapping_unsorted():
    assert merge_intervals([[1, 3], [8, 10], [2, 6], [15, 18]]) == [[1, 6], [8, 10], [15, 18]]


def test_touching_intervals_merge():
    assert merge_intervals([[4, 5], [1, 4]]) == [[1, 5]]


def test_contained_interval():
    assert merge_intervals([[1, 10], [2, 3], [4, 8]]) == [[1, 10]]


def test_disjoint_intervals_sorted_by_start():
    assert merge_intervals([[7, 8], [1, 2], [4, 5]]) == [[1, 2], [4, 5], [7, 8]]


def test_chain_of_overlaps():
    assert merge_intervals([[5, 7], [1, 3], [2, 4], [6, 9]]) == [[1, 4], [5, 9]]


def test_single_and_empty():
    assert merge_intervals([[2, 2]]) == [[2, 2]]
    assert merge_intervals([]) == []


def test_input_not_modified():
    intervals = [[3, 6], [1, 4]]
    assert merge_intervals(intervals) == [[1, 6]]
    assert intervals == [[3, 6], [1, 4]]
