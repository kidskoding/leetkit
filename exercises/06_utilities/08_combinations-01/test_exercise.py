import ast
import inspect

from starter import pairs_with_sum


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_pairs_in_combinations_order():
    assert pairs_with_sum([3, 5, 2, 6, 4], 8) == [(3, 5), (2, 6)]


def test_values_keep_input_order_within_pair():
    assert pairs_with_sum([6, 1, 2, 5], 7) == [(6, 1), (2, 5)]


def test_duplicate_values_form_separate_pairs():
    assert pairs_with_sum([4, 4, 4], 8) == [(4, 4), (4, 4), (4, 4)]


def test_element_not_paired_with_itself():
    assert pairs_with_sum([4, 1], 8) == []


def test_no_match_and_short_input():
    assert pairs_with_sum([1, 2, 3], 100) == []
    assert pairs_with_sum([8], 8) == []
    assert pairs_with_sum([], 0) == []


def test_negative_numbers():
    assert pairs_with_sum([-3, 7, 3, 0], 0) == [(-3, 3)]


def test_uses_combinations():
    assert uses_name(pairs_with_sum, "combinations"), (
        "This drill practices itertools.combinations: loop over combinations(prices, 2)."
    )
