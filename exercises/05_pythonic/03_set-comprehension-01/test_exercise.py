import ast
import inspect

from starter import email_domains


def uses_node(func, node_type):
    tree = ast.parse(inspect.getsource(func))
    return any(isinstance(node, node_type) for node in ast.walk(tree))


def test_basic_domains():
    assert email_domains(["ann@gmail.com", "cy@uiuc.edu"]) == {"gmail.com", "uiuc.edu"}


def test_case_is_normalized():
    assert email_domains(["ann@Gmail.com", "bob@GMAIL.COM"]) == {"gmail.com"}


def test_duplicates_collapse():
    assert email_domains(["a@x.org", "b@x.org", "c@x.org"]) == {"x.org"}


def test_empty_list():
    assert email_domains([]) == set()


def test_uses_last_at_sign():
    assert email_domains(["weird@name@corp.io"]) == {"corp.io"}


def test_returns_a_set():
    assert isinstance(email_domains(["a@b.com"]), set)


def test_uses_set_comprehension():
    assert uses_node(email_domains, ast.SetComp), (
        "This drill practices set comprehensions: return {... for email in emails}-style code."
    )
