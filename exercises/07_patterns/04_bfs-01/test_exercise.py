import copy

from starter import count_islands


def to_grid(rows):
    return [list(row) for row in rows]


def test_classic_three_islands():
    grid = to_grid(["11000", "11000", "00100", "00011"])
    assert count_islands(grid) == 3


def test_one_large_island():
    grid = to_grid(["11110", "11010", "11000", "00000"])
    assert count_islands(grid) == 1


def test_diagonals_do_not_connect():
    grid = to_grid(["101", "010", "101"])
    assert count_islands(grid) == 5


def test_all_water():
    assert count_islands(to_grid(["00", "00"])) == 0


def test_empty_grid():
    assert count_islands([]) == 0
    assert count_islands([[]]) == 0


def test_single_row_and_column():
    assert count_islands(to_grid(["10101"])) == 3
    assert count_islands(to_grid(["1", "1", "0", "1"])) == 2


def test_winding_island_and_grid_unchanged():
    grid = to_grid(["111", "001", "111", "100", "111"])
    original = copy.deepcopy(grid)
    assert count_islands(grid) == 1
    assert grid == original


def test_horizontal_neighbours_connect():
    assert count_islands(to_grid(["0110", "0000", "1101"])) == 3


def test_search_must_also_go_up_and_left():
    # Reaching every cell of these islands from their first cell in reading
    # order needs a step up (the U) or a step left (the L).
    assert count_islands(to_grid(["101", "111"])) == 1
    assert count_islands(to_grid(["01", "11"])) == 1
