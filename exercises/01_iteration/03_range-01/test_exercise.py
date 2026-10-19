import ast
import inspect

from starter import sample_from_end


def calls_builtin(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == name
        for node in ast.walk(tree)
    )


def test_step_reaches_first_element():
    assert sample_from_end([1, 2, 3, 4, 5, 6, 7], 3) == [7, 4, 1]


def test_step_skips_first_element():
    assert sample_from_end([10, 20, 30, 40], 2) == [40, 20]


def test_step_of_one_reverses():
    assert sample_from_end([1, 2, 3], 1) == [3, 2, 1]


def test_step_larger_than_list():
    assert sample_from_end([5, 6, 7], 10) == [7]


def test_single_element():
    assert sample_from_end([42], 2) == [42]


def test_empty_list():
    assert sample_from_end([], 3) == []


def test_uses_range():
    assert calls_builtin(sample_from_end, "range"), (
        "This drill practices range(): generate indices with `range(len(readings) - 1, -1, -k)`."
    )
