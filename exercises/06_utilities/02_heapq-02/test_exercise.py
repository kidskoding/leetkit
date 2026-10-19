from starter import process_tasks


def test_orders_by_priority():
    assert process_tasks([("email", 2), ("deploy", 1), ("lunch", 3)]) == ["deploy", "email", "lunch"]


def test_ties_keep_arrival_order_not_alphabetical():
    assert process_tasks([("b", 1), ("a", 1)]) == ["b", "a"]


def test_mixed_ties():
    tasks = [("review", 2), ("hotfix", 0), ("standup", 2), ("page", 0), ("docs", 5)]
    assert process_tasks(tasks) == ["hotfix", "page", "review", "standup", "docs"]


def test_empty_list():
    assert process_tasks([]) == []


def test_single_task():
    assert process_tasks([("only", 7)]) == ["only"]


def test_negative_priorities_are_most_urgent():
    assert process_tasks([("normal", 0), ("critical", -5), ("low", 3)]) == ["critical", "normal", "low"]
