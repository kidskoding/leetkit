# deque(maxlen=k)

A deque created with `maxlen` never grows past that size. Appending to a full
deque silently drops the item at the opposite end.

```python
from collections import deque

window = deque(maxlen=3)
for x in [1, 2, 3, 4]:
    window.append(x)
print(window)  # deque([2, 3, 4], maxlen=3)
```

The manual version needs an explicit check on every step:

```python
window.append(x)
if len(window) > k:
    window.popleft()
```

## Useful variations

- Keep a running `total` (add the new value, subtract the one about to be evicted) for O(1) averages instead of `sum(window)`.
- `appendleft` on a bounded deque drops from the right end instead.
- `deque(iterable, maxlen=k)` keeps only the last `k` items of any iterable, like a `tail` command.

## Interview relevance

Fixed-size windows show up in stream design questions (Moving Average from
Data Stream, Hit Counter, Logger Rate Limiter) and in "last k elements" style
problems. A bounded deque makes the eviction rule impossible to forget and
signals that you know the standard library well.
