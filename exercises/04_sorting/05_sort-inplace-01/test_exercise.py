import ast
import inspect

from starter import top_k


def uses_keyword(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call) and any(kw.arg == name for kw in node.keywords)
        for node in ast.walk(tree)
    )


def test_largest_first():
    assert top_k([40, 10, 90, 70], 2) == [90, 70]


def test_input_unchanged():
    scores = [40, 10, 90, 70]
    result = top_k(scores, 3)
    assert result == [90, 70, 40]
    assert scores == [40, 10, 90, 70]


def test_duplicates_count_separately():
    assert top_k([3, 8, 8, 1], 3) == [8, 8, 3]


def test_k_larger_than_list():
    assert top_k([5, 1, 5], 10) == [5, 5, 1]


def test_k_zero_and_empty():
    assert top_k([4, 2], 0) == []
    assert top_k([], 3) == []


def test_negative_values():
    assert top_k([-3, -1, -7], 2) == [-1, -3]


def test_uses_reverse():
    assert uses_keyword(top_k, "reverse"), (
        "This drill practices descending sorts: use `sorted(scores, reverse=True)`."
    )
