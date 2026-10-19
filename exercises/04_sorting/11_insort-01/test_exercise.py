import ast
import inspect

from starter import running_medians


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_mixed_stream():
    assert running_medians([5, 15, 1, 3]) == [5.0, 10.0, 5.0, 4.0]


def test_even_count_averages_middle_pair():
    assert running_medians([1, 2]) == [1.0, 1.5]


def test_results_are_floats():
    result = running_medians([4, 8, 6])
    assert result == [4.0, 6.0, 6.0]
    assert all(isinstance(m, float) for m in result)


def test_duplicates_and_negatives():
    assert running_medians([2, 2, 2]) == [2.0, 2.0, 2.0]
    assert running_medians([-2, -8, 4]) == [-2.0, -5.0, -2.0]


def test_descending_stream():
    assert running_medians([9, 7, 5, 3, 1]) == [9.0, 8.0, 7.0, 6.0, 5.0]


def test_empty_and_input_not_modified():
    assert running_medians([]) == []
    stream = [3, 1, 2]
    running_medians(stream)
    assert stream == [3, 1, 2]


def test_uses_insort():
    assert uses_name(running_medians, "insort"), (
        "This drill practices bisect.insort: add each value with `insort(window, value)` instead of sorting."
    )
