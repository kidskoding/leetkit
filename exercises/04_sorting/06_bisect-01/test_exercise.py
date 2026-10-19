import ast
import inspect

from starter import search_insert


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_target_found():
    assert search_insert([1, 3, 5, 6], 5) == 2


def test_insert_in_middle():
    assert search_insert([1, 3, 5, 6], 2) == 1


def test_insert_at_end():
    assert search_insert([1, 3, 5, 6], 7) == 4


def test_insert_at_start():
    assert search_insert([1, 3, 5, 6], 0) == 0


def test_negative_values():
    assert search_insert([-10, -4, 0, 9], -5) == 1
    assert search_insert([-10, -4, 0, 9], -10) == 0


def test_empty_and_single():
    assert search_insert([], 4) == 0
    assert search_insert([4], 4) == 0
    assert search_insert([4], 5) == 1


def test_uses_bisect_left():
    assert uses_name(search_insert, "bisect_left"), (
        "This drill practices bisect: return `bisect_left(nums, target)`."
    )
