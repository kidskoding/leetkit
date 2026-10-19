from starter import climb_stairs


def test_small_staircases():
    assert climb_stairs(2) == 2
    assert climb_stairs(3) == 3


def test_base_cases():
    assert climb_stairs(0) == 1
    assert climb_stairs(1) == 1


def test_medium_staircase():
    assert climb_stairs(5) == 8
    assert climb_stairs(10) == 89


def test_twenty_five_steps():
    assert climb_stairs(25) == 121393


def test_is_memoized_with_functools():
    assert hasattr(climb_stairs, "cache_info"), (
        "Decorate climb_stairs with @cache or @lru_cache so repeated subproblems are reused."
    )
