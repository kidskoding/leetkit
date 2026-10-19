import ast
import inspect

from starter import count_correct


def calls_builtin(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == name
        for node in ast.walk(tree)
    )


def test_some_correct():
    assert count_correct(["a", "c", "b", "d"], ["a", "b", "b", "d"]) == 3


def test_all_correct():
    assert count_correct([1, 2, 3], [1, 2, 3]) == 3


def test_none_correct():
    assert count_correct(["a", "b"], ["b", "a"]) == 0


def test_shorter_answers_ignores_extra_key():
    assert count_correct(["a", "b"], ["a", "b", "c"]) == 2


def test_longer_answers_ignores_extra_answers():
    assert count_correct(["a", "x", "c", "d"], ["a", "b"]) == 1


def test_empty_lists():
    assert count_correct([], []) == 0
    assert count_correct([], ["a"]) == 0


def test_uses_zip():
    assert calls_builtin(count_correct, "zip"), (
        "This drill practices zip(): iterate with `for given, correct in zip(answers, key)`."
    )
