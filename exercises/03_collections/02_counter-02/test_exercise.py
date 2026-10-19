import ast
import inspect

from starter import top_k_frequent


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_basic_top_two():
    assert top_k_frequent([1, 1, 1, 2, 2, 3], 2) == [1, 2]


def test_most_frequent_first():
    assert top_k_frequent([3, 1, 3, 2, 3, 2], 3) == [3, 2, 1]


def test_ties_keep_first_seen_order():
    assert top_k_frequent([4, 5, 5, 4, 6], 2) == [4, 5]


def test_k_larger_than_distinct_values():
    assert top_k_frequent([7, 7, 8], 5) == [7, 8]


def test_k_zero():
    assert top_k_frequent([1, 2, 3], 0) == []


def test_empty_list():
    assert top_k_frequent([], 3) == []


def test_negative_values():
    assert top_k_frequent([-1, -1, 0, 2, 2, 2], 1) == [2]


def test_uses_most_common():
    assert uses_name(top_k_frequent, "most_common"), (
        "This drill practices Counter.most_common(k): call it instead of sorting counts by hand."
    )
