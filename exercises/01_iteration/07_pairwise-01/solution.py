def largest_swing(temps):
    return max((abs(b - a) for a, b in zip(temps, temps[1:])), default=0)
