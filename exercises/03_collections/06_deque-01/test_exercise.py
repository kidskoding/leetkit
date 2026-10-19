from starter import min_flights


def test_shortcut_is_used():
    assert min_flights(4, [(0, 1), (1, 2), (2, 3), (0, 2)], 0, 3) == 2


def test_direct_flight():
    assert min_flights(2, [(0, 1)], 0, 1) == 1


def test_routes_work_both_ways():
    assert min_flights(3, [(0, 1), (1, 2)], 2, 0) == 2


def test_unreachable_city():
    assert min_flights(3, [(0, 1)], 0, 2) == -1


def test_same_city():
    assert min_flights(1, [], 0, 0) == 0


def test_chooses_shorter_of_two_routes():
    flights = [(0, 1), (1, 2), (2, 3), (3, 4), (0, 5), (5, 4)]
    assert min_flights(6, flights, 0, 4) == 2


def test_graph_with_cycle():
    flights = [(0, 1), (1, 2), (2, 0), (2, 3)]
    assert min_flights(4, flights, 0, 3) == 2
