from starter import course_order


def test_diamond_dependencies():
    assert course_order(4, [(1, 0), (2, 0), (3, 1), (3, 2)]) == [0, 1, 2, 3]


def test_no_prerequisites_is_ascending():
    assert course_order(3, []) == [0, 1, 2]


def test_unlocked_smaller_course_jumps_ahead():
    assert course_order(4, [(1, 0)]) == [0, 1, 2, 3]
    assert course_order(5, [(0, 4), (1, 3)]) == [2, 3, 1, 4, 0]


def test_two_course_cycle():
    assert course_order(2, [(0, 1), (1, 0)]) == []


def test_cycle_elsewhere_still_fails():
    assert course_order(4, [(1, 0), (2, 1), (1, 2)]) == []


def test_self_prerequisite_is_a_cycle():
    assert course_order(2, [(1, 1)]) == []


def test_duplicate_pairs_and_zero_courses():
    assert course_order(2, [(1, 0), (1, 0)]) == [0, 1]
    assert course_order(0, []) == []
