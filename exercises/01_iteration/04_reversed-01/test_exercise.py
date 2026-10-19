import ast
import inspect

from starter import last_index


def calls_builtin(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == name
        for node in ast.walk(tree)
    )


def test_last_of_several():
    assert last_index(["ana", "bo", "ana", "cy"], "ana") == 2


def test_single_occurrence():
    assert last_index(["ana", "bo", "cy"], "bo") == 1


def test_match_at_ends():
    assert last_index(["cy", "bo", "ana"], "ana") == 2
    assert last_index(["cy", "bo", "ana"], "cy") == 0


def test_all_same():
    assert last_index(["bo", "bo", "bo"], "bo") == 2


def test_missing_user():
    assert last_index(["ana", "bo"], "dee") == -1


def test_empty_list():
    assert last_index([], "ana") == -1


def test_uses_reversed():
    assert calls_builtin(last_index, "reversed"), (
        "This drill practices reversed(): loop with `for i in reversed(range(len(events)))`."
    )
