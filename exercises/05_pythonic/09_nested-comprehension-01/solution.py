def make_grid(rows, cols, fill):
    return [[fill for _ in range(cols)] for _ in range(rows)]


def flatten(matrix):
    return [x for row in matrix for x in row]
