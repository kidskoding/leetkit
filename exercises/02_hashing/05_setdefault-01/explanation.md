# dict.setdefault()

`d.setdefault(key, default)` returns `d[key]` if the key exists. If it does not,
it first stores `default` under `key` and then returns it. Because it returns
the stored object itself, you can mutate it in place:

```python
groups = {}
for customer, order_id in orders:
    groups.setdefault(customer, []).append(order_id)
```

The branching version does the same in four lines:

```python
if customer not in groups:
    groups[customer] = []
groups[customer].append(order_id)
```

Since dicts keep insertion order, the keys come out in the order each customer
was first seen, and each list grows in input order.

## Useful variations

- `defaultdict(list)` from `collections` lets you write `groups[key].append(x)` directly.
- `d.setdefault(key, set()).add(x)` groups into sets when you only need unique values.
- `d.setdefault(key, 0)` just reads with a default and stores it; for counting, `d.get(key, 0) + 1` is clearer.
- The default is built on every call, even when the key exists, so keep it cheap (`[]`, not an expensive call).

## Interview relevance

Grouping by key is the first step of many graph and data problems: building an
adjacency list from `edges`, collecting a user's transactions before checking a
fraud rule, or bucketing logs by service. `setdefault` keeps that step to one
line so the interviewer can focus on what you do with the groups.
