from starter import LRUCache


def test_get_missing_key():
    cache = LRUCache(2)
    assert cache.get(1) == -1


def test_put_then_get():
    cache = LRUCache(2)
    cache.put(1, 100)
    cache.put(2, 200)
    assert cache.get(1) == 100
    assert cache.get(2) == 200


def test_evicts_least_recently_used():
    cache = LRUCache(2)
    cache.put(1, 1)
    cache.put(2, 2)
    cache.put(3, 3)
    assert cache.get(1) == -1
    assert cache.get(2) == 2
    assert cache.get(3) == 3


def test_get_refreshes_recency():
    cache = LRUCache(2)
    cache.put(1, 1)
    cache.put(2, 2)
    assert cache.get(1) == 1
    cache.put(3, 3)
    assert cache.get(2) == -1
    assert cache.get(1) == 1


def test_update_refreshes_without_evicting():
    cache = LRUCache(2)
    cache.put(1, 1)
    cache.put(2, 2)
    cache.put(1, 10)
    assert cache.get(2) == 2
    assert cache.get(1) == 10
    cache.put(3, 3)
    assert cache.get(2) == -1
    assert cache.get(1) == 10


def test_capacity_one():
    cache = LRUCache(1)
    cache.put("a", 1)
    cache.put("b", 2)
    assert cache.get("a") == -1
    assert cache.get("b") == 2


def test_leetcode_sequence():
    cache = LRUCache(2)
    cache.put(1, 1)
    cache.put(2, 2)
    assert cache.get(1) == 1
    cache.put(3, 3)
    assert cache.get(2) == -1
    cache.put(4, 4)
    assert cache.get(1) == -1
    assert cache.get(3) == 3
    assert cache.get(4) == 4
