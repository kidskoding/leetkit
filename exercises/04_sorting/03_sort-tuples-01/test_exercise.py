import ast
import inspect

from starter import rank_players


def key_returns_tuple(func):
    tree = ast.parse(inspect.getsource(func))
    for node in ast.walk(tree):
        if isinstance(node, ast.Call):
            for kw in node.keywords:
                if kw.arg == "key" and any(isinstance(n, ast.Tuple) for n in ast.walk(kw.value)):
                    return True
    return False


def test_highest_score_first():
    assert rank_players([("ana", 70), ("bo", 95), ("cy", 80)]) == ["bo", "cy", "ana"]


def test_ties_broken_by_name_ascending():
    assert rank_players([("ana", 90), ("bo", 95), ("cy", 90)]) == ["bo", "ana", "cy"]


def test_tie_breaker_not_reversed():
    assert rank_players([("zed", 10), ("amy", 10), ("max", 10)]) == ["amy", "max", "zed"]


def test_negative_scores():
    assert rank_players([("a", -5), ("b", 0), ("c", -1)]) == ["b", "c", "a"]


def test_single_and_empty():
    assert rank_players([("solo", 1)]) == ["solo"]
    assert rank_players([]) == []


def test_uses_tuple_key():
    assert key_returns_tuple(rank_players), (
        "This drill practices tuple keys: sort with `key=lambda p: (-p[1], p[0])` in one call."
    )
