import ast
import inspect

from starter import split_route


def uses_starred_assignment(func):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Assign)
        and isinstance(node.targets[0], (ast.Tuple, ast.List))
        and any(isinstance(elt, ast.Starred) for elt in node.targets[0].elts)
        for node in ast.walk(tree)
    )


def test_two_stops():
    assert split_route("LAX-ORD-DEN-JFK") == ("LAX", ["ORD", "DEN"], "JFK")


def test_one_stop():
    assert split_route("BOS-ATL-MIA") == ("BOS", ["ATL"], "MIA")


def test_direct_flight_has_empty_stops():
    assert split_route("SFO-SEA") == ("SFO", [], "SEA")


def test_stops_is_a_list():
    _, stops, _ = split_route("A-B-C-D-E")
    assert stops == ["B", "C", "D"]
    assert isinstance(stops, list)


def test_round_trip_route():
    assert split_route("ORD-DEN-ORD") == ("ORD", ["DEN"], "ORD")


def test_uses_starred_unpacking():
    assert uses_starred_assignment(split_route), (
        "This drill practices extended unpacking: write `origin, *stops, destination = route.split('-')`."
    )
