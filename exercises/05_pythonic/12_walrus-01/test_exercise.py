import ast
import inspect

from starter import high_scoring_words


def uses_node(func, node_type):
    tree = ast.parse(inspect.getsource(func))
    return any(isinstance(node, node_type) for node in ast.walk(tree))


def test_basic_filter():
    assert high_scoring_words(["cab", "zoo", "a"], 5) == [("cab", 6), ("zoo", 56)]


def test_score_equal_to_cutoff_is_excluded():
    assert high_scoring_words(["bad", "ace"], 9) == []


def test_empty_list():
    assert high_scoring_words([], 0) == []


def test_keeps_input_order_and_duplicates():
    assert high_scoring_words(["zz", "b", "zz"], 1) == [("zz", 52), ("b", 2), ("zz", 52)]


def test_empty_word_scores_zero():
    assert high_scoring_words(["", "a"], -1) == [("", 0), ("a", 1)]


def test_all_words_qualify():
    assert high_scoring_words(["hello", "world"], 0) == [("hello", 52), ("world", 72)]


def test_uses_walrus():
    assert uses_node(high_scoring_words, ast.NamedExpr), (
        "This drill practices the walrus operator: write `if (score := ...) > cutoff` in the comprehension."
    )
