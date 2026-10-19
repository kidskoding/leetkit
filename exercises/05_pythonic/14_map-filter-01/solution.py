def positive_readings(tokens):
    return list(filter(lambda n: n > 0, map(int, tokens)))
