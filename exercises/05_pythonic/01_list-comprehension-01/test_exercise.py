import ast
import inspect

from starter import even_squares


def uses_node(func, node_type):
    tree = ast.parse(inspect.getsource(func))
    return any(isinstance(node, node_type) for node in ast.walk(tree))


def test_mixed_numbers():
    assert even_squares([1, 2, 3, 4]) == [4, 16]


def test_keeps_input_order():
    assert even_squares([6, 1, 2, 8]) == [36, 4, 64]


def test_no_even_numbers():
    assert even_squares([1, 3, 5]) == []


def test_empty_list():
    assert even_squares([]) == []


def test_zero_and_negatives():
    assert even_squares([0, -2, -3, -4]) == [0, 4, 16]


def test_duplicates_are_kept():
    assert even_squares([2, 2, 3]) == [4, 4]


def test_uses_list_comprehension():
    assert uses_node(even_squares, ast.ListComp), (
        "This drill practices list comprehensions: return [n * n for n in nums if n % 2 == 0]-style code."
    )
