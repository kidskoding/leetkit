from starter import longest_unique_run


def test_classic_example():
    assert longest_unique_run("abcabcbb") == 3


def test_all_same_character():
    assert longest_unique_run("bbbbb") == 1


def test_answer_not_a_prefix():
    assert longest_unique_run("pwwkew") == 3


def test_empty_string():
    assert longest_unique_run("") == 0


def test_stale_index_outside_window_is_ignored():
    assert longest_unique_run("abba") == 2
    assert longest_unique_run("tmmzuxt") == 5


def test_all_unique():
    assert longest_unique_run("abcdef") == 6


def test_spaces_and_case_count():
    assert longest_unique_run("a b!a") == 4
    assert longest_unique_run("aA") == 2
