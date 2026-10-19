def rotate_clockwise(matrix):
    return [list(row) for row in zip(*reversed(matrix))]
