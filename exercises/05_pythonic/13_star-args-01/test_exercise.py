import ast
import inspect

from starter import merge_settings


def has_star_params(func):
    tree = ast.parse(inspect.getsource(func))
    definition = tree.body[0]
    return definition.args.vararg is not None and definition.args.kwarg is not None


def test_later_layer_wins():
    assert merge_settings({"theme": "light", "font": 12}, {"font": 14}) == {"theme": "light", "font": 14}


def test_keyword_override_beats_layers():
    assert merge_settings({"font": 12}, {"font": 14}, font=16, debug=True) == {"font": 16, "debug": True}


def test_no_arguments():
    assert merge_settings() == {}


def test_unpacked_list_of_layers():
    configs = [{"a": 1}, {"b": 2}, {"a": 3, "c": 4}]
    assert merge_settings(*configs) == {"a": 3, "b": 2, "c": 4}


def test_unpacked_keyword_dict():
    flags = {"verbose": True, "retries": 5}
    assert merge_settings({"retries": 3}, **flags) == {"retries": 5, "verbose": True}


def test_inputs_are_not_modified():
    base = {"x": 1}
    extra = {"x": 2}
    result = merge_settings(base, extra, y=3)
    assert result == {"x": 2, "y": 3}
    assert base == {"x": 1}
    assert extra == {"x": 2}
    assert result is not base


def test_signature_uses_star_args_and_kwargs():
    assert has_star_params(merge_settings), (
        "This drill practices variadic parameters: define merge_settings(*layers, **overrides)."
    )
