import ast
import inspect

from starter import over_budget


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_some_over():
    assert over_budget({"rent": 1200, "food": 450, "fun": 300}, 400) == ["rent", "food"]


def test_keeps_insertion_order():
    assert over_budget({"zoo": 9, "apple": 8, "mango": 1}, 5) == ["zoo", "apple"]


def test_equal_to_limit_is_not_over():
    assert over_budget({"food": 400, "gas": 401}, 400) == ["gas"]


def test_none_over():
    assert over_budget({"food": 10, "gas": 20}, 100) == []


def test_negative_limit():
    assert over_budget({"refund": -5, "fee": 0}, -1) == ["fee"]


def test_empty_dict():
    assert over_budget({}, 0) == []


def test_uses_items():
    assert uses_name(over_budget, "items"), (
        "This drill practices dict.items(): loop with `for category, amount in spending.items()`."
    )
