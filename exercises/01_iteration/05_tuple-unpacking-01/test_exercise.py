import ast
import inspect

from starter import fib


def uses_tuple_assignment(func):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Assign) and isinstance(node.targets[0], ast.Tuple)
        for node in ast.walk(tree)
    )


def test_base_cases():
    assert fib(0) == 0
    assert fib(1) == 1


def test_small_values():
    assert [fib(n) for n in range(2, 8)] == [1, 2, 3, 5, 8, 13]


def test_ten():
    assert fib(10) == 55


def test_large_value():
    assert fib(90) == 2880067194370816120


def test_uses_tuple_assignment():
    assert uses_tuple_assignment(fib), (
        "This drill practices tuple unpacking: update both values at once with `a, b = b, a + b`."
    )
