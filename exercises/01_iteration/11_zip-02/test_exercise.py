import ast
import inspect

import pytest

from starter import price_lookup


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_basic_pairs():
    assert price_lookup(["apple", "bread"], [3, 5]) == {"apple": 3, "bread": 5}


def test_repeated_item_keeps_last_price():
    assert price_lookup(["tea", "milk", "tea"], [2, 1, 4]) == {"tea": 4, "milk": 1}


def test_single_item():
    assert price_lookup(["salt"], [0]) == {"salt": 0}


def test_empty_lists():
    assert price_lookup([], []) == {}


def test_more_items_than_prices_raises():
    with pytest.raises(ValueError):
        price_lookup(["apple", "bread", "eggs"], [3, 5])


def test_more_prices_than_items_raises():
    with pytest.raises(ValueError):
        price_lookup(["apple"], [3, 5])


def test_uses_zip():
    assert uses_name(price_lookup, "zip"), (
        "This drill practices zip(): build the dict with `dict(zip(items, prices, strict=True))`."
    )
