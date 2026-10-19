import ast
import inspect

from starter import points_in_box


def loops_with_nested_unpacking(func):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, (ast.For, ast.comprehension))
        and isinstance(node.target, ast.Tuple)
        and any(isinstance(elt, ast.Tuple) for elt in node.target.elts)
        for node in ast.walk(tree)
    )


BOX = ((0, 0), (3, 3))


def test_some_inside():
    points = [((1, 1), "cafe"), ((5, 2), "park"), ((2, 3), "gym")]
    assert points_in_box(points, BOX) == ["cafe", "gym"]


def test_edges_and_corners_count():
    points = [((0, 0), "a"), ((3, 3), "b"), ((3, 0), "c"), ((0, 2), "d")]
    assert points_in_box(points, BOX) == ["a", "b", "c", "d"]


def test_outside_on_each_side():
    points = [((-1, 1), "left"), ((4, 1), "right"), ((1, -1), "below"), ((1, 4), "above")]
    assert points_in_box(points, BOX) == []


def test_negative_coordinates_box():
    points = [((-2, -2), "in"), ((-6, 0), "out"), ((0, 0), "corner")]
    assert points_in_box(points, ((-5, -5), (0, 0))) == ["in", "corner"]


def test_keeps_input_order_and_duplicates():
    points = [((2, 2), "z"), ((1, 1), "a"), ((2, 2), "z")]
    assert points_in_box(points, BOX) == ["z", "a", "z"]


def test_empty_points():
    assert points_in_box([], BOX) == []


def test_uses_nested_unpacking():
    assert loops_with_nested_unpacking(points_in_box), (
        "This drill practices nested unpacking: loop with `for (x, y), label in points` "
        "(in a for statement or a comprehension) instead of indexing entry[0][0]."
    )
