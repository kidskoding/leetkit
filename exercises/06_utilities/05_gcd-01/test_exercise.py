import ast
import inspect

from starter import bundle_size


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_common_divisor():
    assert bundle_size([12, 18, 30]) == 6


def test_coprime_counts():
    assert bundle_size([7, 9]) == 1


def test_single_count():
    assert bundle_size([42]) == 42


def test_zero_counts_are_ignored():
    assert bundle_size([0, 15, 25]) == 5


def test_empty_and_all_zero():
    assert bundle_size([]) == 0
    assert bundle_size([0, 0]) == 0


def test_uses_gcd():
    assert uses_name(bundle_size, "gcd"), (
        "This drill practices math.gcd: fold it across the list, e.g. reduce(math.gcd, counts, 0)."
    )
