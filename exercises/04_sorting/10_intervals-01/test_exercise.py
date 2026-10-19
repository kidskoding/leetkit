from starter import can_attend_all


def test_overlap_after_sorting():
    assert can_attend_all([[0, 30], [5, 10], [15, 20]]) is False


def test_no_overlap_unsorted_input():
    assert can_attend_all([[7, 10], [2, 4]]) is True


def test_touching_meetings_are_allowed():
    assert can_attend_all([[7, 10], [2, 4], [4, 7]]) is True


def test_overlap_only_between_late_neighbours():
    assert can_attend_all([[1, 2], [3, 4], [9, 12], [5, 6], [11, 13]]) is False


def test_same_start_time_overlaps():
    assert can_attend_all([[3, 5], [3, 4]]) is False


def test_empty_and_single():
    assert can_attend_all([]) is True
    assert can_attend_all([[5, 8]]) is True


def test_input_not_modified():
    meetings = [[9, 10], [1, 2]]
    can_attend_all(meetings)
    assert meetings == [[9, 10], [1, 2]]
