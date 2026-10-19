import ast
import inspect

import pytest

from starter import moving_average


def uses_bounded_deque(func):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call)
        and isinstance(node.func, ast.Name)
        and node.func.id == "deque"
        and any(kw.arg == "maxlen" for kw in node.keywords)
        for node in ast.walk(tree)
    )


def test_window_of_two():
    assert moving_average([1, 2, 3, 4], 2) == pytest.approx([1.0, 1.5, 2.5, 3.5])


def test_window_of_three():
    assert moving_average([3, 6, 9, 12, 15], 3) == pytest.approx([3.0, 4.5, 6.0, 9.0, 12.0])


def test_k_larger_than_input():
    assert moving_average([10, 20, 30], 5) == pytest.approx([10.0, 15.0, 20.0])


def test_window_of_one():
    assert moving_average([4, -2, 7], 1) == pytest.approx([4.0, -2.0, 7.0])


def test_non_terminating_decimals():
    assert moving_average([1, 1, 2, 5], 3) == pytest.approx([1.0, 1.0, 4 / 3, 8 / 3])


def test_empty_readings():
    assert moving_average([], 3) == []


def test_uses_deque_maxlen():
    assert uses_bounded_deque(moving_average), (
        "This drill practices bounded deques: create the window with `deque(maxlen=k)`."
    )
