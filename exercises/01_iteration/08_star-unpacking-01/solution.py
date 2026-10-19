def split_route(route):
    origin, *stops, destination = route.split("-")
    return origin, stops, destination
