import ast
import inspect

from starter import closest_point


def calls_with_key(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call)
        and isinstance(node.func, ast.Name)
        and node.func.id == name
        and any(kw.arg == "key" for kw in node.keywords)
        for node in ast.walk(tree)
    )


def test_closest_point():
    assert closest_point([(3, 4), (1, -1), (5, 0)]) == (1, -1)


def test_negative_coordinates():
    assert closest_point([(-6, -8), (-2, 1), (4, 4)]) == (-2, 1)


def test_tie_returns_first():
    assert closest_point([(0, 2), (2, 0), (-2, 0)]) == (0, 2)


def test_origin_itself():
    assert closest_point([(1, 1), (0, 0), (0, 0)]) == (0, 0)


def test_single_point():
    assert closest_point([(7, -3)]) == (7, -3)


def test_empty_list():
    assert closest_point([]) is None


def test_uses_min_with_key():
    assert calls_with_key(closest_point, "min"), (
        "This drill practices min() with key=: `min(points, key=lambda p: ..., default=None)`."
    )
