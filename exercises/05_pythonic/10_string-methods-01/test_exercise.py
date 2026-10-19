import ast
import inspect

from starter import is_valid_palindrome


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_sentence_with_punctuation():
    assert is_valid_palindrome("A man, a plan, a canal: Panama") is True


def test_not_a_palindrome():
    assert is_valid_palindrome("race a car") is False


def test_empty_string():
    assert is_valid_palindrome("") is True


def test_only_punctuation():
    assert is_valid_palindrome(" .,!") is True


def test_case_insensitive():
    assert is_valid_palindrome("Aa") is True


def test_digits_count():
    assert is_valid_palindrome("0P") is False
    assert is_valid_palindrome("1b-2-B1") is True


def test_uses_isalnum():
    assert uses_name(is_valid_palindrome, "isalnum"), (
        "This drill practices str methods: filter characters with ch.isalnum() before comparing."
    )
