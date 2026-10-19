import ast
import inspect

from starter import group_by_letters


def uses_name(func, name):
    tree = ast.parse(inspect.getsource(func))
    return any(
        (isinstance(node, ast.Name) and node.id == name)
        or (isinstance(node, ast.Attribute) and node.attr == name)
        for node in ast.walk(tree)
    )


def test_groups_by_letter_set():
    assert group_by_letters(["loop", "stop", "pool", "lop", "post"]) == [
        ["loop", "pool", "lop"],
        ["stop", "post"],
    ]


def test_repeated_letters_do_not_matter():
    assert group_by_letters(["aab", "ab", "abc"]) == [["aab", "ab"], ["abc"]]


def test_empty_list():
    assert group_by_letters([]) == []


def test_repeated_words_are_kept():
    assert group_by_letters(["tab", "bat", "tab"]) == [["tab", "bat", "tab"]]


def test_case_sensitive():
    assert group_by_letters(["Ab", "ab", "ba"]) == [["Ab"], ["ab", "ba"]]


def test_empty_string_is_its_own_group():
    assert group_by_letters(["", "a", ""]) == [["", ""], ["a"]]


def test_uses_frozenset():
    assert uses_name(group_by_letters, "frozenset"), (
        "This drill practices frozenset keys: group with `groups.setdefault(frozenset(word), []).append(word)`."
    )
