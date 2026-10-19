import ast
import inspect

from starter import merge_streams


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_three_streams():
    assert merge_streams([[1, 4, 9], [2, 3, 10], [5]]) == [1, 2, 3, 4, 5, 9, 10]


def test_duplicates_are_kept():
    assert merge_streams([[3, 3], [], [3]]) == [3, 3, 3]
    assert merge_streams([[1, 2, 2], [2, 5]]) == [1, 2, 2, 2, 5]


def test_single_stream_unchanged():
    assert merge_streams([[0, 7, 8]]) == [0, 7, 8]


def test_no_streams_or_all_empty():
    assert merge_streams([]) == []
    assert merge_streams([[], []]) == []


def test_negative_timestamps():
    assert merge_streams([[-5, 0], [-7, -1, 3]]) == [-7, -5, -1, 0, 3]


def test_one_stream_entirely_after_another():
    assert merge_streams([[10, 11], [1, 2]]) == [1, 2, 10, 11]


def test_uses_heapq_merge():
    assert uses_name(merge_streams, "merge"), (
        "This drill practices heapq.merge: return list(heapq.merge(*streams))."
    )
