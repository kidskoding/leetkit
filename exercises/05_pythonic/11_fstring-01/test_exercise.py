import ast
import inspect

from starter import receipt_line


def uses_format_spec(func):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.FormattedValue) and node.format_spec is not None
        for node in ast.walk(tree)
    )


def test_basic_line():
    assert receipt_line("apple", 3, 0.5) == "apple        3     1.50"


def test_total_is_rounded_to_two_places():
    assert receipt_line("pen", 3, 0.333) == "pen          3     1.00"


def test_zero_quantity():
    assert receipt_line("gum", 0, 1.25) == "gum          0     0.00"


def test_item_exactly_ten_characters():
    assert receipt_line("toothpaste", 1, 3.0) == "toothpaste   1     3.00"


def test_long_item_is_not_truncated():
    assert receipt_line("watermelon slice", 12, 2.25) == "watermelon slice  12    27.00"


def test_wide_total_is_not_truncated():
    assert receipt_line("laptop", 100, 1234.5) == "laptop     100 123450.00"


def test_uses_fstring_format_spec():
    assert uses_format_spec(receipt_line), (
        "This drill practices f-string format specs: write something like f\"{total:>8.2f}\"."
    )
