import ast
import inspect

from starter import run_length_encode


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_mixed_runs():
    assert run_length_encode("aaabcc") == "a3b1c2"


def test_no_repeats():
    assert run_length_encode("abc") == "a1b1c1"


def test_empty_string():
    assert run_length_encode("") == ""


def test_single_character():
    assert run_length_encode("z") == "z1"


def test_long_run_uses_multi_digit_count():
    assert run_length_encode("x" * 12) == "x12"


def test_character_can_start_a_new_run():
    assert run_length_encode("aabaa") == "a2b1a2"


def test_case_sensitive():
    assert run_length_encode("aA") == "a1A1"


def test_uses_join():
    assert uses_name(run_length_encode, "join"), (
        "This drill practices str.join(): append pieces to a list and return \"\".join(parts)."
    )
