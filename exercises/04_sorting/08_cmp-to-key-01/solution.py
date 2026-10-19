from functools import cmp_to_key


def largest_number(nums):
    def compare(a, b):
        if a + b > b + a:
            return -1
        if a + b < b + a:
            return 1
        return 0

    pieces = sorted(map(str, nums), key=cmp_to_key(compare))
    return "0" if pieces[0] == "0" else "".join(pieces)
