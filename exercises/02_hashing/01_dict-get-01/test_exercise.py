import ast
import inspect

from starter import char_counts


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_repeated_letters():
    assert char_counts("hello") == {"h": 1, "e": 1, "l": 2, "o": 1}


def test_empty_string():
    assert char_counts("") == {}


def test_single_character():
    assert char_counts("z") == {"z": 1}


def test_all_same_character():
    assert char_counts("aaaa") == {"a": 4}


def test_case_sensitive_and_spaces():
    assert char_counts("Aa a!") == {"A": 1, "a": 2, " ": 1, "!": 1}


def test_uses_dict_get():
    assert uses_name(char_counts, "get"), (
        "This drill practices dict.get(): update with `counts[ch] = counts.get(ch, 0) + 1`."
    )
