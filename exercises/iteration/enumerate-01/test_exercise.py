import ast
import inspect

from starter import find_index


def calls_builtin(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == name
        for node in ast.walk(tree)
    )


def test_target_present():
    assert find_index([10, 20, 30], 20) == 1


def test_first_occurrence():
    assert find_index([7, 3, 7, 3], 3) == 1


def test_target_at_start_and_end():
    assert find_index([5, 6, 7], 5) == 0
    assert find_index([5, 6, 7], 7) == 2


def test_target_missing():
    assert find_index([10, 20, 30], 50) == -1


def test_empty_list():
    assert find_index([], 5) == -1


def test_uses_enumerate():
    assert calls_builtin(find_index, "enumerate"), (
        "This drill practices enumerate(): loop with `for i, num in enumerate(nums)`."
    )
