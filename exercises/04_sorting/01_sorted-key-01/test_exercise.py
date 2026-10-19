import ast
import inspect

from starter import sort_by_length


def calls_with_key(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call)
        and (
            (isinstance(node.func, ast.Name) and node.func.id == name)
            or (isinstance(node.func, ast.Attribute) and node.func.attr == name)
        )
        and any(kw.arg == "key" for kw in node.keywords)
        for node in ast.walk(tree)
    )


def test_orders_by_length():
    assert sort_by_length(["banana", "fig", "apple", "kiwi"]) == ["fig", "kiwi", "apple", "banana"]


def test_ties_keep_input_order():
    assert sort_by_length(["bb", "aa", "c", "dd"]) == ["c", "bb", "aa", "dd"]


def test_not_alphabetical():
    assert sort_by_length(["zz", "a", "yy", "b"]) == ["a", "b", "zz", "yy"]


def test_empty_list():
    assert sort_by_length([]) == []


def test_single_word():
    assert sort_by_length(["solo"]) == ["solo"]


def test_input_not_modified():
    words = ["ccc", "a", "bb"]
    result = sort_by_length(words)
    assert result == ["a", "bb", "ccc"]
    assert words == ["ccc", "a", "bb"]


def test_uses_sorted_with_key():
    assert calls_with_key(sort_by_length, "sorted"), (
        "This drill practices sorted() with a key function: `sorted(words, key=len)`."
    )
