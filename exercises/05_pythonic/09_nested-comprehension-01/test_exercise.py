import ast
import inspect

from starter import flatten, make_grid


def comprehensions(func):
    tree = ast.parse(inspect.getsource(func))
    return [node for node in ast.walk(tree) if isinstance(node, ast.ListComp)]


def test_make_grid_shape_and_fill():
    assert make_grid(2, 3, 0) == [[0, 0, 0], [0, 0, 0]]


def test_make_grid_rows_are_independent():
    grid = make_grid(3, 2, ".")
    grid[0][0] = "#"
    assert grid == [["#", "."], [".", "."], [".", "."]]


def test_make_grid_zero_rows_or_cols():
    assert make_grid(0, 3, 0) == []
    assert make_grid(2, 0, 0) == [[], []]


def test_flatten_rows_in_order():
    assert flatten([[1, 2, 3], [4, 5, 6]]) == [1, 2, 3, 4, 5, 6]


def test_flatten_ragged_and_empty_rows():
    assert flatten([[1, 2], [], [3]]) == [1, 2, 3]


def test_flatten_empty_matrix():
    assert flatten([]) == []


def test_uses_nested_comprehensions():
    grid_comps = comprehensions(make_grid)
    nested = any(isinstance(comp.elt, ast.ListComp) for comp in grid_comps)
    two_clause = any(len(comp.generators) == 2 for comp in comprehensions(flatten))
    assert nested and two_clause, (
        "This drill practices nested comprehensions: make_grid should be "
        "[[fill for _ in range(cols)] for _ in range(rows)] and flatten should be "
        "[x for row in matrix for x in row]."
    )
