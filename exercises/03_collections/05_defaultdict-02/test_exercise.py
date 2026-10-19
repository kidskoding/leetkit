import ast
import inspect

from starter import build_graph


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_path_graph():
    assert build_graph([(1, 2), (2, 3)]) == {1: [2], 2: [1, 3], 3: [2]}


def test_neighbours_sorted():
    assert build_graph([(5, 9), (5, 1), (5, 3)]) == {5: [1, 3, 9], 9: [5], 1: [5], 3: [5]}


def test_duplicate_edges_collapse():
    assert build_graph([(1, 2), (2, 1), (1, 2)]) == {1: [2], 2: [1]}


def test_disconnected_components():
    assert build_graph([(0, 1), (7, 8)]) == {0: [1], 1: [0], 7: [8], 8: [7]}


def test_empty_edges():
    assert build_graph([]) == {}


def test_string_names():
    edges = [("ana", "bo"), ("cy", "ana")]
    assert build_graph(edges) == {"ana": ["bo", "cy"], "bo": ["ana"], "cy": ["ana"]}


def test_uses_defaultdict():
    assert uses_name(build_graph, "defaultdict"), (
        "This drill practices defaultdict: build the graph with `graph = defaultdict(set)`."
    )
