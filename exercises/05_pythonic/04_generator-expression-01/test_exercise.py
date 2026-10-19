import ast
import inspect

from starter import cart_total


def uses_node(func, node_type):
    tree = ast.parse(inspect.getsource(func))
    return any(isinstance(node, node_type) for node in ast.walk(tree))


def test_two_items():
    assert cart_total([(250, 2), (100, 3)]) == 800


def test_single_item():
    assert cart_total([(999, 1)]) == 999


def test_empty_cart():
    assert cart_total([]) == 0


def test_zero_quantity_contributes_nothing():
    assert cart_total([(500, 0), (200, 2)]) == 400


def test_many_items():
    assert cart_total([(1, 1), (2, 2), (3, 3), (4, 4)]) == 30


def test_uses_generator_expression():
    assert uses_node(cart_total, ast.GeneratorExp), (
        "This drill practices generator expressions: write sum(price * qty for price, qty in items), "
        "with no square brackets."
    )
