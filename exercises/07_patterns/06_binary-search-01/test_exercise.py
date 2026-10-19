from starter import first_bad_build


def test_bad_in_middle():
    assert first_bad_build(5, lambda b: b >= 4) == 4


def test_first_build_is_bad():
    assert first_bad_build(10, lambda b: True) == 1


def test_only_last_build_is_bad():
    assert first_bad_build(10, lambda b: b >= 10) == 10


def test_no_bad_build():
    assert first_bad_build(3, lambda b: False) == -1


def test_single_and_zero_builds():
    assert first_bad_build(1, lambda b: True) == 1
    assert first_bad_build(1, lambda b: False) == -1
    assert first_bad_build(0, lambda b: True) == -1


def test_uses_logarithmic_checks_within_range():
    n = 1_000_000
    calls = []

    def is_bad(b):
        calls.append(b)
        return b >= 765_432

    assert first_bad_build(n, is_bad) == 765_432
    assert all(1 <= b <= n for b in calls), "Only call is_bad on builds 1..n."
    assert len(calls) <= 25, "Use binary search: about log2(n) calls to is_bad, not one per build."


def test_no_bad_build_stays_within_range():
    calls = []

    def is_bad(b):
        calls.append(b)
        return False

    assert first_bad_build(100, is_bad) == -1
    assert all(1 <= b <= 100 for b in calls), "Only call is_bad on builds 1..n."
