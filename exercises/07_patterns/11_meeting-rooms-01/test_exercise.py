from starter import min_meeting_rooms


def test_long_meeting_overlaps_others():
    assert min_meeting_rooms([(0, 30), (5, 10), (15, 20)]) == 2


def test_back_to_back_shares_a_room():
    assert min_meeting_rooms([(1, 5), (5, 10)]) == 1


def test_all_overlap():
    assert min_meeting_rooms([(1, 10), (2, 9), (3, 8)]) == 3


def test_unsorted_input():
    assert min_meeting_rooms([(15, 20), (0, 30), (10, 15), (5, 10)]) == 2


def test_rooms_freed_and_reused():
    assert min_meeting_rooms([(1, 4), (2, 5), (4, 7), (5, 8), (9, 10)]) == 2


def test_empty_and_single():
    assert min_meeting_rooms([]) == 0
    assert min_meeting_rooms([(3, 4)]) == 1


def test_identical_meetings():
    assert min_meeting_rooms([(2, 6), (2, 6), (2, 6), (6, 9)]) == 3
