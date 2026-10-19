import ast
import inspect

from starter import range_sums


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_several_queries():
    assert range_sums([2, 4, 1, 3], [(0, 1), (1, 3), (2, 2)]) == [6, 8, 1]


def test_whole_array():
    assert range_sums([1, 2, 3, 4, 5], [(0, 4)]) == [15]


def test_single_element_ranges():
    assert range_sums([9, 8, 7], [(0, 0), (2, 2)]) == [9, 7]


def test_negative_numbers():
    assert range_sums([-2, 0, 3, -5, 2, -1], [(0, 2), (2, 5), (0, 5)]) == [1, -1, -3]


def test_no_queries():
    assert range_sums([1, 2, 3], []) == []


def test_uses_accumulate():
    assert uses_name(range_sums, "accumulate"), (
        "This drill practices itertools.accumulate: prefix = list(accumulate(nums, initial=0))."
    )
