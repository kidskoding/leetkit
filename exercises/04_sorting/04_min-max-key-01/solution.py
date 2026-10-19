def closest_point(points):
    return min(points, key=lambda p: p[0] * p[0] + p[1] * p[1], default=None)
