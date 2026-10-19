import ast
import inspect

from starter import is_strong_password


def calls_builtin(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == name
        for node in ast.walk(tree)
    )


def test_strong_password():
    assert is_strong_password("Secure123") is True


def test_too_short():
    assert is_strong_password("Ab1cdef") is False


def test_missing_uppercase():
    assert is_strong_password("secure123") is False


def test_missing_lowercase():
    assert is_strong_password("SECURE123") is False


def test_missing_digit():
    assert is_strong_password("SecurePass") is False


def test_contains_whitespace():
    assert is_strong_password("Secure 123") is False
    assert is_strong_password("Secure123\t") is False


def test_empty_string():
    assert is_strong_password("") is False


def test_uses_any_and_all():
    assert calls_builtin(is_strong_password, "any") and calls_builtin(is_strong_password, "all"), (
        "This drill practices any() and all(): use any(...) for the 'at least one' rules "
        "and all(...) for the 'no whitespace' rule."
    )


def test_exactly_eight_characters_is_long_enough():
    assert is_strong_password("Abcdef12") is True
