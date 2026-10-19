import ast
import inspect

from starter import word_lengths


def uses_node(func, node_type):
    tree = ast.parse(inspect.getsource(func))
    return any(isinstance(node, node_type) for node in ast.walk(tree))


def test_basic_words():
    assert word_lengths(["apple", "fig"]) == {"apple": 5, "fig": 3}


def test_empty_list():
    assert word_lengths([]) == {}


def test_duplicates_appear_once():
    assert word_lengths(["kiwi", "kiwi", "pear"]) == {"kiwi": 4, "pear": 4}


def test_case_sensitive_keys():
    assert word_lengths(["Go", "go"]) == {"Go": 2, "go": 2}


def test_empty_string_word():
    assert word_lengths(["", "a"]) == {"": 0, "a": 1}


def test_uses_dict_comprehension():
    assert uses_node(word_lengths, ast.DictComp), (
        "This drill practices dict comprehensions: return {word: len(word) for word in words}-style code."
    )
