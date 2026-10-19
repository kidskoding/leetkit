import ast
import inspect

from starter import positive_readings


def calls_function(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == name
        for node in ast.walk(tree)
    )


def test_mixed_tokens():
    assert positive_readings(["42", "-7", " 3 ", "0"]) == [42, 3]


def test_empty_list():
    assert positive_readings([]) == []


def test_all_non_positive():
    assert positive_readings(["0", "-1", "-20"]) == []


def test_keeps_order_and_duplicates():
    assert positive_readings(["5", "1", "5", "-5"]) == [5, 1, 5]


def test_plus_sign_and_whitespace():
    assert positive_readings(["+8", "\t9\n", "  -0 "]) == [8, 9]


def test_returns_a_list():
    assert isinstance(positive_readings(["1"]), list)


def test_uses_map_and_filter():
    assert calls_function(positive_readings, "map") and calls_function(positive_readings, "filter"), (
        "This drill practices map() and filter(): convert with map(int, tokens) and select with filter(...)."
    )
