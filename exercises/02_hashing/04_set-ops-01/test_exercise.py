import ast
import inspect

from starter import compare_rosters


def uses_set_operators(func):
    tree = ast.parse(inspect.getsource(func))
    ops = {type(node.op) for node in ast.walk(tree) if isinstance(node, ast.BinOp)}
    return {ast.BitAnd, ast.BitOr, ast.Sub, ast.BitXor} <= ops


def test_partial_overlap():
    assert compare_rosters(["ana", "ben", "cal"], ["ben", "dee"]) == {
        "both": ["ben"],
        "either": ["ana", "ben", "cal", "dee"],
        "only_a": ["ana", "cal"],
        "exactly_one": ["ana", "cal", "dee"],
    }


def test_no_overlap():
    assert compare_rosters(["zed", "amy"], ["kim"]) == {
        "both": [],
        "either": ["amy", "kim", "zed"],
        "only_a": ["amy", "zed"],
        "exactly_one": ["amy", "kim", "zed"],
    }


def test_identical_rosters():
    assert compare_rosters(["bo", "al"], ["al", "bo"]) == {
        "both": ["al", "bo"],
        "either": ["al", "bo"],
        "only_a": [],
        "exactly_one": [],
    }


def test_repeats_are_collapsed():
    assert compare_rosters(["ana", "ana", "ben"], ["ben", "ben"]) == {
        "both": ["ben"],
        "either": ["ana", "ben"],
        "only_a": ["ana"],
        "exactly_one": ["ana"],
    }


def test_empty_inputs():
    empty = {"both": [], "either": [], "only_a": [], "exactly_one": []}
    assert compare_rosters([], []) == empty
    assert compare_rosters([], ["kim"]) == {
        "both": [],
        "either": ["kim"],
        "only_a": [],
        "exactly_one": ["kim"],
    }


def test_uses_set_operators():
    assert uses_set_operators(compare_rosters), (
        "This drill practices set operators: compute the results with a & b, a | b, a - b and a ^ b."
    )
