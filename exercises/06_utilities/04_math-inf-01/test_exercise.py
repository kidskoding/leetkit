from starter import max_profit


def test_buy_low_sell_high():
    assert max_profit([7, 1, 5, 3, 6, 4]) == 5


def test_falling_prices_give_zero():
    assert max_profit([7, 6, 4, 3, 1]) == 0


def test_minimum_after_maximum():
    assert max_profit([2, 9, 1, 4]) == 7


def test_empty_and_single_price():
    assert max_profit([]) == 0
    assert max_profit([5]) == 0


def test_flat_prices():
    assert max_profit([3, 3, 3]) == 0


def test_large_prices():
    assert max_profit([10**12, 10**12 + 7]) == 7
