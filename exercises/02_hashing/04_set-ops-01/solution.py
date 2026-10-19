def compare_rosters(course_a, course_b):
    a, b = set(course_a), set(course_b)
    return {
        "both": sorted(a & b),
        "either": sorted(a | b),
        "only_a": sorted(a - b),
        "exactly_one": sorted(a ^ b),
    }
