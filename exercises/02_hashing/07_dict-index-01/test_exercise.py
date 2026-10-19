from starter import first_unique_index


def test_first_character_unique():
    assert first_unique_index("leetcode") == 0


def test_unique_in_middle():
    assert first_unique_index("loveleetcode") == 2


def test_no_unique_character():
    assert first_unique_index("aabb") == -1


def test_empty_string():
    assert first_unique_index("") == -1


def test_single_character():
    assert first_unique_index("x") == 0


def test_unique_at_end():
    assert first_unique_index("aabbc") == 4


def test_case_sensitive():
    assert first_unique_index("aAa") == 1


def test_character_seen_three_times_is_not_unique():
    assert first_unique_index("abaca") == 1
