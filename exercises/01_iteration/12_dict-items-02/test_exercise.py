from starter import top_sellers


def test_top_two():
    sales = {"mug": 40, "pen": 95, "cap": 60, "bag": 12}
    assert top_sellers(sales, 2) == [("pen", 95), ("cap", 60)]


def test_ties_broken_by_name():
    sales = {"tea": 30, "gum": 30, "ink": 50}
    assert top_sellers(sales, 3) == [("ink", 50), ("gum", 30), ("tea", 30)]


def test_tie_break_ignores_insertion_order():
    sales = {"zip": 5, "bat": 5, "map": 5}
    assert top_sellers(sales, 2) == [("bat", 5), ("map", 5)]


def test_n_larger_than_catalog_returns_all():
    assert top_sellers({"a": 1, "b": 2}, 10) == [("b", 2), ("a", 1)]


def test_n_zero():
    assert top_sellers({"a": 1, "b": 2}, 0) == []


def test_empty_sales():
    assert top_sellers({}, 3) == []


def test_zero_units_ranked_last():
    assert top_sellers({"dud": 0, "hit": 7, "ok": 3}, 3) == [("hit", 7), ("ok", 3), ("dud", 0)]


def test_input_not_modified():
    sales = {"mug": 40, "pen": 95}
    top_sellers(sales, 1)
    assert sales == {"mug": 40, "pen": 95}
    assert list(sales) == ["mug", "pen"]
