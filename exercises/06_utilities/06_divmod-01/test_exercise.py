import ast
import inspect

from starter import format_duration


def calls_builtin(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == name
        for node in ast.walk(tree)
    )


def test_hours_minutes_seconds():
    assert format_duration(3725) == "1:02:05"


def test_zero():
    assert format_duration(0) == "0:00:00"


def test_under_a_minute():
    assert format_duration(59) == "0:00:59"


def test_exact_boundaries():
    assert format_duration(60) == "0:01:00"
    assert format_duration(3600) == "1:00:00"
    assert format_duration(3599) == "0:59:59"


def test_many_hours_are_not_wrapped():
    assert format_duration(100 * 3600 + 7) == "100:00:07"


def test_uses_divmod():
    assert calls_builtin(format_duration, "divmod"), (
        "This drill practices divmod(): minutes, secs = divmod(seconds, 60)."
    )
