from starter import find_peaks


def test_single_center_peak():
    assert find_peaks([[1, 2, 1], [2, 5, 2], [1, 2, 1]]) == [(1, 1)]


def test_corner_and_edge_peaks_have_fewer_neighbours():
    grid = [
        [9, 1, 7],
        [1, 0, 1],
    ]
    assert find_peaks(grid) == [(0, 0), (0, 2)]


def test_plateau_is_not_a_peak():
    assert find_peaks([[3, 3], [1, 0]]) == []


def test_diagonals_do_not_count():
    grid = [
        [9, 1],
        [1, 8],
    ]
    assert find_peaks(grid) == [(0, 0), (1, 1)]


def test_row_major_order():
    grid = [
        [1, 0, 0, 4],
        [0, 0, 0, 0],
        [5, 0, 2, 0],
    ]
    assert find_peaks(grid) == [(0, 0), (0, 3), (2, 0), (2, 2)]


def test_negative_heights_single_row():
    assert find_peaks([[-5, -1, -3, -4]]) == [(0, 1)]


def test_one_by_one_grid():
    assert find_peaks([[42]]) == [(0, 0)]


def test_empty_grid():
    assert find_peaks([]) == []


def test_taller_cell_directly_above_or_below_blocks_a_peak():
    assert find_peaks([[9, 0], [5, 1]]) == [(0, 0)]
    assert find_peaks([[5, 1], [9, 0]]) == [(1, 0)]
