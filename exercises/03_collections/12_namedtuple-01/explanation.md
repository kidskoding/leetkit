# collections.namedtuple

`namedtuple` builds a small tuple subclass whose positions also have names.
Records stay immutable, hashable, and unpackable, but read like objects.

```python
from collections import namedtuple

Player = namedtuple("Player", ["name", "score"])
p = Player("ana", 90)
p.score          # 90
name, score = p  # still unpacks like a tuple
p == ("ana", 90) # True
```

The plain-tuple version works, but `p[1]` says nothing about which field it
is, and swapping two positions by mistake is easy to miss.

## Useful variations

- `p._replace(score=95)` returns a copy with one field changed.
- `Player._make(["ana", 90])` builds a record from any iterable; `p._asdict()` turns it into a dict.
- `typing.NamedTuple` gives the same record with a class syntax and type hints.

## Interview relevance

Graph and scheduling problems pass around small records like edges, events,
and heap entries. Naming the fields makes your code self-explaining while you
talk through it, and namedtuples still sort and hash, so they work as dict
keys, set members, and heap items.
