import ast
import inspect

from starter import rotate_right


def uses_node(func, node_type):
    tree = ast.parse(inspect.getsource(func))
    return any(isinstance(node, node_type) for node in ast.walk(tree))


def test_basic_rotation():
    assert rotate_right([1, 2, 3, 4, 5], 2) == [4, 5, 1, 2, 3]


def test_k_zero_returns_same_order():
    assert rotate_right([1, 2, 3], 0) == [1, 2, 3]


def test_k_equal_to_length():
    assert rotate_right([1, 2, 3], 3) == [1, 2, 3]


def test_k_larger_than_length():
    assert rotate_right([1, 2, 3], 4) == [3, 1, 2]
    assert rotate_right([1, 2, 3], 10) == [3, 1, 2]


def test_empty_list():
    assert rotate_right([], 3) == []


def test_single_element():
    assert rotate_right([7], 5) == [7]


def test_input_is_not_modified():
    nums = [1, 2, 3, 4]
    rotate_right(nums, 1)
    assert nums == [1, 2, 3, 4]


def test_uses_slicing():
    assert uses_node(rotate_right, ast.Slice), (
        "This drill practices slicing: build the result from slices such as nums[-k:] and nums[:-k]."
    )
