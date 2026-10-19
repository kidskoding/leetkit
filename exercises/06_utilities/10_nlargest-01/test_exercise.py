import ast
import inspect

from starter import top_orders


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


ORDERS = [
    {"id": "a1", "total": 25},
    {"id": "b2", "total": 90},
    {"id": "c3", "total": 40},
    {"id": "d4", "total": 12.5},
]


def test_top_two():
    assert top_orders(ORDERS, 2) == ["b2", "c3"]


def test_largest_first_for_all():
    assert top_orders(ORDERS, 4) == ["b2", "c3", "a1", "d4"]


def test_k_larger_than_list():
    assert top_orders(ORDERS[:2], 5) == ["b2", "a1"]


def test_ties_keep_input_order():
    orders = [{"id": "x", "total": 50}, {"id": "y", "total": 50}, {"id": "z", "total": 70}]
    assert top_orders(orders, 2) == ["z", "x"]
    assert top_orders(orders, 3) == ["z", "x", "y"]


def test_zero_k_and_empty_orders():
    assert top_orders(ORDERS, 0) == []
    assert top_orders([], 3) == []


def test_uses_nlargest():
    assert uses_name(top_orders, "nlargest"), (
        "This drill practices heapq.nlargest: call heapq.nlargest(k, orders, key=...)."
    )
