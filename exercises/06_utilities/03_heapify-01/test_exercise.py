import ast
import inspect

from starter import k_smallest


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_two_smallest():
    assert k_smallest([40, 10, 30, 20], 2) == [10, 20]


def test_duplicates_count_separately():
    assert k_smallest([3, 1, 3, 1, 2], 3) == [1, 1, 2]


def test_k_larger_than_list():
    assert k_smallest([5, 1], 5) == [1, 5]


def test_k_zero_and_empty_list():
    assert k_smallest([4, 2], 0) == []
    assert k_smallest([], 3) == []


def test_input_is_not_mutated():
    prices = [9, 4, 7, 1, 8]
    assert k_smallest(prices, 2) == [1, 4]
    assert prices == [9, 4, 7, 1, 8]


def test_uses_heapify():
    assert uses_name(k_smallest, "heapify"), (
        "This drill practices heapq.heapify(): copy the list, then call heapq.heapify(copy)."
    )
