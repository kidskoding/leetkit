from starter import rotate_clockwise


def test_two_by_two():
    assert rotate_clockwise([[1, 2], [3, 4]]) == [[3, 1], [4, 2]]


def test_three_by_three():
    grid = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
    assert rotate_clockwise(grid) == [[7, 4, 1], [8, 5, 2], [9, 6, 3]]


def test_rectangular_grid_changes_shape():
    assert rotate_clockwise([[1, 2, 3], [4, 5, 6]]) == [[4, 1], [5, 2], [6, 3]]


def test_single_row_becomes_column():
    assert rotate_clockwise([[1, 2, 3]]) == [[1], [2], [3]]


def test_rows_are_lists():
    result = rotate_clockwise([[1, 2], [3, 4]])
    assert all(isinstance(row, list) for row in result)


def test_input_not_modified():
    grid = [[1, 2], [3, 4]]
    rotate_clockwise(grid)
    assert grid == [[1, 2], [3, 4]]


def test_four_rotations_restore_grid():
    grid = [[1, 2, 3], [4, 5, 6]]
    result = grid
    for _ in range(4):
        result = rotate_clockwise(result)
    assert result == grid


def test_empty_matrix():
    assert rotate_clockwise([]) == []
    assert rotate_clockwise([[]]) == []
