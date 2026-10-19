import ast
import inspect

from starter import orders_by_customer


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_groups_orders():
    assert orders_by_customer([("ava", 101), ("raj", 102), ("ava", 103)]) == {
        "ava": [101, 103],
        "raj": [102],
    }


def test_empty_list():
    assert orders_by_customer([]) == {}


def test_single_order():
    assert orders_by_customer([("kim", 7)]) == {"kim": [7]}


def test_keys_in_first_seen_order():
    result = orders_by_customer([("zoe", 1), ("ava", 2), ("zoe", 3), ("max", 4)])
    assert list(result) == ["zoe", "ava", "max"]


def test_lists_keep_input_order_and_repeats():
    assert orders_by_customer([("ava", 30), ("ava", 10), ("ava", 30)]) == {"ava": [30, 10, 30]}


def test_customer_names_are_case_sensitive():
    assert orders_by_customer([("Ava", 1), ("ava", 2)]) == {"Ava": [1], "ava": [2]}


def test_uses_setdefault():
    assert uses_name(orders_by_customer, "setdefault"), (
        "This drill practices dict.setdefault(): group with `groups.setdefault(customer, []).append(order_id)`."
    )
