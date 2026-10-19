from starter import count_clusters


def test_two_clusters():
    assert count_clusters(5, [(0, 1), (1, 2), (3, 4)]) == 2


def test_no_links_every_server_alone():
    assert count_clusters(3, []) == 3


def test_zero_servers():
    assert count_clusters(0, []) == 0


def test_isolated_server_counts():
    assert count_clusters(4, [(0, 1), (1, 2)]) == 2


def test_cycle_is_one_cluster():
    assert count_clusters(4, [(0, 1), (1, 2), (2, 3), (3, 0)]) == 1


def test_duplicate_links_and_self_loops():
    assert count_clusters(4, [(0, 1), (1, 0), (0, 1), (2, 2)]) == 3


def test_links_listed_out_of_order():
    assert count_clusters(6, [(5, 3), (0, 4), (3, 1), (4, 2)]) == 2
