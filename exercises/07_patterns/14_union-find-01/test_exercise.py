from starter import find_redundant_cable


def test_triangle_with_string_labels():
    assert find_redundant_cable([("a", "b"), ("b", "c"), ("a", "c")]) == ("a", "c")


def test_no_redundant_cable():
    assert find_redundant_cable([(1, 2), (3, 4)]) is None
    assert find_redundant_cable([]) is None


def test_first_redundant_in_input_order():
    cables = [(1, 2), (2, 3), (3, 4), (1, 4), (1, 5), (5, 3)]
    assert find_redundant_cable(cables) == (1, 4)


def test_groups_merge_before_loop_closes():
    assert find_redundant_cable([(1, 2), (3, 4), (2, 3), (4, 1)]) == (4, 1)


def test_self_loop_is_redundant():
    assert find_redundant_cable([("x", "y"), ("z", "z")]) == ("z", "z")


def test_reversed_duplicate_cable():
    assert find_redundant_cable([(7, 8), (8, 7)]) == (8, 7)


def test_long_chain():
    cables = [(i, i + 1) for i in range(3000)] + [(0, 3000), (5, 6)]
    assert find_redundant_cable(cables) == (0, 3000)
