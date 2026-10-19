import ast
import inspect

from starter import number_lines


def enumerate_has_start(func):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(node, ast.Call)
        and isinstance(node.func, ast.Name)
        and node.func.id == "enumerate"
        and (len(node.args) == 2 or any(kw.arg == "start" for kw in node.keywords))
        for node in ast.walk(tree)
    )


def test_two_lines():
    assert number_lines(["revenue up", "costs flat"]) == ["1. revenue up", "2. costs flat"]


def test_single_line():
    assert number_lines(["only"]) == ["1. only"]


def test_blank_line_still_numbered():
    assert number_lines(["intro", "", "end"]) == ["1. intro", "2. ", "3. end"]


def test_double_digit_numbers():
    result = number_lines([f"row {i}" for i in range(12)])
    assert result[0] == "1. row 0"
    assert result[9] == "10. row 9"
    assert result[11] == "12. row 11"


def test_duplicate_lines_get_their_own_numbers():
    assert number_lines(["ok", "ok"]) == ["1. ok", "2. ok"]


def test_empty_list():
    assert number_lines([]) == []


def test_uses_enumerate_with_start():
    assert enumerate_has_start(number_lines), (
        "This drill practices enumerate(..., start=1): loop with `for n, line in enumerate(lines, start=1)`."
    )
