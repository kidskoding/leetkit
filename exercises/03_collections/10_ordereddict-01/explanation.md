# OrderedDict as an LRU list

An `OrderedDict` remembers insertion order and can reorder keys in O(1). That
makes it a hash map and a recency list in one object: keep the least recently
used key at the front and the newest at the back.

```python
from collections import OrderedDict

d = OrderedDict(a=1, b=2, c=3)
d.move_to_end("a")      # order is now b, c, a
d.popitem(last=False)   # removes ('b', 2), the front item
```

Without it you would write a doubly linked list of nodes plus a dict from key
to node, and splice nodes by hand on every get and put.

## Useful variations

- `move_to_end(key, last=False)` moves a key to the front instead.
- `popitem()` with the default `last=True` pops the newest item, which gives a stack-like eviction.
- A plain `dict` also keeps insertion order, but has no `move_to_end`; you would `pop` and reinsert the key.

## Interview relevance

LRU Cache appears constantly in design rounds and as a follow-up to caching
discussions. Using OrderedDict gets you a correct O(1) answer quickly, and
explaining that it is a linked list plus hash map underneath shows you know
what the interviewer is really testing.
