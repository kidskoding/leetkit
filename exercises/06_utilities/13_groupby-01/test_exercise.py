import ast
import inspect

from starter import compress_levels


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_runs_in_order():
    assert compress_levels(["INFO", "INFO", "WARN", "ERROR", "ERROR", "ERROR"]) == [
        ("INFO", 2),
        ("WARN", 1),
        ("ERROR", 3),
    ]


def test_repeat_after_gap_starts_new_run():
    assert compress_levels(["INFO", "INFO", "WARN", "INFO"]) == [("INFO", 2), ("WARN", 1), ("INFO", 1)]


def test_alternating_levels():
    assert compress_levels(["A", "B", "A", "B"]) == [("A", 1), ("B", 1), ("A", 1), ("B", 1)]


def test_single_run():
    assert compress_levels(["DEBUG"] * 5) == [("DEBUG", 5)]


def test_single_entry():
    assert compress_levels(["ERROR"]) == [("ERROR", 1)]


def test_empty():
    assert compress_levels([]) == []


def test_uses_groupby():
    assert uses_name(compress_levels, "groupby"), (
        "This drill practices itertools.groupby: loop with `for level, run in groupby(levels)`."
    )
