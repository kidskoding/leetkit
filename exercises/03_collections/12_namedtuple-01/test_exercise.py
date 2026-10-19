import ast
import inspect
import sys

from starter import best_player, parse_players


def uses_name(module, name):
    tree = ast.parse(inspect.getsource(module))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_records_have_named_fields():
    player = parse_players(["ana,90"])[0]
    assert player.name == "ana"
    assert player.score == 90
    assert player._fields == ("name", "score")


def test_records_are_tuples_in_input_order():
    assert parse_players(["ana,90", "bo,75"]) == [("ana", 90), ("bo", 75)]


def test_negative_score_parsed_as_int():
    assert parse_players(["zed,-5"])[0].score == -5


def test_best_player():
    assert best_player(["ana,90", "bo,75", "cy,88"]) == "ana"


def test_tie_goes_to_first():
    assert best_player(["ana,90", "bo,95", "cy,95"]) == "bo"


def test_empty_input():
    assert parse_players([]) == []
    assert best_player([]) is None


def test_uses_namedtuple():
    assert uses_name(sys.modules[parse_players.__module__], "namedtuple"), (
        "This drill practices namedtuple: define `Player = namedtuple(\"Player\", [\"name\", \"score\"])`."
    )
