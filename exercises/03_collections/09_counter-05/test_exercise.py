import ast
import inspect

from starter import common_chars


def uses_intersection(func):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, (ast.BinOp, ast.AugAssign)) and isinstance(node.op, ast.BitAnd)
        for node in ast.walk(tree)
    )


def test_keeps_duplicates():
    assert common_chars(["bella", "label", "roller"]) == ["e", "l", "l"]


def test_minimum_count_wins():
    assert common_chars(["aaa", "aa", "aaaa"]) == ["a", "a"]


def test_result_is_sorted():
    assert common_chars(["cool", "lock", "cook"]) == ["c", "o"]


def test_nothing_in_common():
    assert common_chars(["abc", "def"]) == []


def test_single_word_returns_its_letters_sorted():
    assert common_chars(["banana"]) == ["a", "a", "a", "b", "n", "n"]


def test_empty_word_clears_everything():
    assert common_chars(["abc", "", "cab"]) == []


def test_uses_counter_intersection():
    assert uses_intersection(common_chars), (
        "This drill practices Counter intersection: combine counts with `a & b` or `common &= Counter(word)`."
    )
