def points_in_box(points, box):
    (x_min, y_min), (x_max, y_max) = box
    return [label for (x, y), label in points if x_min <= x <= x_max and y_min <= y <= y_max]
