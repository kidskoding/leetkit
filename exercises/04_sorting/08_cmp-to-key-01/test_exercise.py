from starter import largest_number


def test_two_numbers():
    assert largest_number([10, 2]) == "210"


def test_shared_prefixes():
    assert largest_number([3, 30, 34, 5, 9]) == "9534330"


def test_prefix_that_should_go_second():
    assert largest_number([121, 12]) == "12121"
    assert largest_number([432, 43243]) == "43243432"


def test_all_zeros_collapse_to_single_zero():
    assert largest_number([0, 0]) == "0"
    assert largest_number([0, 0, 0, 0]) == "0"


def test_zeros_mixed_with_other_numbers():
    assert largest_number([0, 5, 0]) == "500"


def test_single_number():
    assert largest_number([7]) == "7"
    assert largest_number([0]) == "0"


def test_large_values():
    assert largest_number([999999998, 999999997, 999999999]) == "999999999999999998999999997"
