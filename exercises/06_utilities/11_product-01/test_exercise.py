import ast
import inspect

from starter import all_codes


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_two_letters_length_two():
    assert all_codes("ab", 2) == ["aa", "ab", "ba", "bb"]


def test_order_follows_alphabet_not_sorted():
    assert all_codes("ba", 2) == ["bb", "ba", "ab", "aa"]


def test_length_one_is_the_alphabet():
    assert all_codes("xyz", 1) == ["x", "y", "z"]


def test_binary_strings_length_three():
    assert all_codes("01", 3) == ["000", "001", "010", "011", "100", "101", "110", "111"]


def test_count_is_alphabet_size_to_the_n():
    codes = all_codes("0123", 4)
    assert len(codes) == 4 ** 4
    assert codes[0] == "0000" and codes[-1] == "3333"


def test_zero_length_and_empty_alphabet():
    assert all_codes("xy", 0) == [""]
    assert all_codes("", 2) == []


def test_uses_product():
    assert uses_name(all_codes, "product"), (
        "This drill practices itertools.product: loop over product(alphabet, repeat=n)."
    )
