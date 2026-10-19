from starter import largest_swing


def test_mixed_changes():
    assert largest_swing([12, 15, 9, 10]) == 6


def test_drop_counts_as_swing():
    assert largest_swing([30, 10, 12]) == 20


def test_negative_readings():
    assert largest_swing([-5, 3, -2]) == 8


def test_flat_readings():
    assert largest_swing([7, 7, 7]) == 0


def test_only_adjacent_pairs_count():
    assert largest_swing([0, 5, 10, 15]) == 5


def test_two_readings():
    assert largest_swing([4, 1]) == 3


def test_fewer_than_two_readings():
    assert largest_swing([20]) == 0
    assert largest_swing([]) == 0
