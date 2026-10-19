from starter import nearest_restaurants


def test_single_closest():
    assert nearest_restaurants([(1, 3), (-2, 2)], 1) == [(-2, 2)]


def test_closest_first_order():
    assert nearest_restaurants([(3, 3), (5, -1), (-2, 4)], 2) == [(3, 3), (-2, 4)]


def test_ties_keep_input_order():
    points = [(0, 5), (1, 0), (0, 1), (-1, 0), (5, 0)]
    assert nearest_restaurants(points, 3) == [(1, 0), (0, 1), (-1, 0)]


def test_k_covers_all_points():
    points = [(2, 2), (0, 0), (-1, 1)]
    assert nearest_restaurants(points, 5) == [(0, 0), (-1, 1), (2, 2)]


def test_k_zero_and_empty_input():
    assert nearest_restaurants([(1, 1)], 0) == []
    assert nearest_restaurants([], 3) == []


def test_negative_coordinates_and_duplicates():
    points = [(-3, -4), (3, 4), (1, -1), (1, -1)]
    assert nearest_restaurants(points, 3) == [(1, -1), (1, -1), (-3, -4)]
