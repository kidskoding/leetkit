import ast
import inspect

from starter import rank_words


def key_returns_tuple(func):
    tree = ast.parse(inspect.getsource(func))
    for node in ast.walk(tree):
        if isinstance(node, ast.Call):
            for kw in node.keywords:
                if kw.arg == "key" and any(isinstance(n, ast.Tuple) for n in ast.walk(kw.value)):
                    return True
    return False


def test_most_frequent_first():
    assert rank_words(["pear", "fig", "pear", "apple", "fig", "pear"]) == ["pear", "fig", "apple"]


def test_ties_broken_alphabetically():
    assert rank_words(["kiwi", "banana", "kiwi", "banana", "date"]) == ["banana", "kiwi", "date"]


def test_all_unique_is_alphabetical():
    assert rank_words(["zeta", "alpha", "mu"]) == ["alpha", "mu", "zeta"]


def test_case_sensitive():
    assert rank_words(["apple", "Apple", "apple"]) == ["apple", "Apple"]
    assert rank_words(["b", "B", "a"]) == ["B", "a", "b"]


def test_single_and_empty():
    assert rank_words(["solo", "solo"]) == ["solo"]
    assert rank_words([]) == []


def test_uses_tuple_key():
    assert key_returns_tuple(rank_words), (
        "This drill practices tuple keys: sort with `key=lambda w: (-counts[w], w)` in one call."
    )
