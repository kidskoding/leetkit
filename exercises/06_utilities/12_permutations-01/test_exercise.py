import ast
import inspect

from starter import arrangements


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_distinct_letters():
    assert arrangements("cat") == ["act", "atc", "cat", "cta", "tac", "tca"]


def test_repeated_letters_are_deduplicated():
    assert arrangements("aab") == ["aab", "aba", "baa"]


def test_all_same_letter():
    assert arrangements("zzz") == ["zzz"]


def test_single_letter():
    assert arrangements("q") == ["q"]


def test_empty_word():
    assert arrangements("") == [""]


def test_count_with_two_pairs():
    result = arrangements("aabb")
    assert len(result) == 6
    assert result[0] == "aabb" and result[-1] == "bbaa"


def test_uses_permutations():
    assert uses_name(arrangements, "permutations"), (
        "This drill practices itertools.permutations: loop over permutations(word)."
    )
