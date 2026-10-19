from starter import is_happy


def test_happy_example():
    assert is_happy(19) is True


def test_unhappy_example():
    assert is_happy(2) is False


def test_one_is_happy():
    assert is_happy(1) is True


def test_powers_of_ten_are_happy():
    assert is_happy(10) is True
    assert is_happy(100) is True


def test_value_on_the_cycle_is_unhappy():
    assert is_happy(4) is False
    assert is_happy(89) is False


def test_matches_known_happy_numbers_below_50():
    happy = {1, 7, 10, 13, 19, 23, 28, 31, 32, 44, 49}
    assert [k for k in range(1, 50) if is_happy(k)] == sorted(happy)


def test_large_number():
    assert is_happy(1111111) is True
