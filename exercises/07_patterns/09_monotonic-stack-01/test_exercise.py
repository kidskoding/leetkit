from starter import days_until_warmer


def test_mixed_week():
    assert days_until_warmer([73, 74, 75, 71, 69, 72, 76, 73]) == [1, 1, 4, 2, 1, 1, 0, 0]


def test_equal_temperature_is_not_warmer():
    assert days_until_warmer([70, 70, 71]) == [2, 1, 0]


def test_strictly_falling_never_warms():
    assert days_until_warmer([90, 80, 70, 60]) == [0, 0, 0, 0]


def test_strictly_rising_waits_one_day():
    assert days_until_warmer([30, 40, 50, 60]) == [1, 1, 1, 0]


def test_one_hot_day_resolves_many():
    assert days_until_warmer([50, 40, 45, 30, 35, 60]) == [5, 1, 3, 1, 1, 0]


def test_empty_and_single_day():
    assert days_until_warmer([]) == []
    assert days_until_warmer([55]) == [0]


def test_negative_temperatures():
    assert days_until_warmer([-5, -10, -3, -3, 0]) == [2, 1, 2, 1, 0]
