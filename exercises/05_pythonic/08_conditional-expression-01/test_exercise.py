import ast
import inspect

from starter import fizzbuzz


def comprehension_uses_ifexp(func):
    tree = ast.parse(inspect.getsource(func))
    return any(
        isinstance(comp, ast.ListComp) and any(isinstance(node, ast.IfExp) for node in ast.walk(comp.elt))
        for comp in ast.walk(tree)
    )


def test_first_five():
    assert fizzbuzz(5) == ["1", "2", "Fizz", "4", "Buzz"]


def test_fifteen_is_fizzbuzz():
    assert fizzbuzz(15)[-1] == "FizzBuzz"


def test_full_fifteen():
    assert fizzbuzz(15) == [
        "1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz",
        "11", "Fizz", "13", "14", "FizzBuzz",
    ]


def test_zero_returns_empty_list():
    assert fizzbuzz(0) == []


def test_one():
    assert fizzbuzz(1) == ["1"]


def test_numbers_are_strings():
    assert all(isinstance(item, str) for item in fizzbuzz(7))


def test_uses_conditional_expression_in_comprehension():
    assert comprehension_uses_ifexp(fizzbuzz), (
        "This drill practices conditional expressions: return "
        "[\"Fizz\" if ... else ... for i in range(1, n + 1)]-style code."
    )
