# Two dicts for a one-to-one mapping

A mapping is one-to-one when every key has a single value and every value comes
from a single key. A dict only enforces the first half, so keep a second dict
for the reverse direction:

```python
s_to_t = {}
t_to_s = {}
for a, b in zip(s, t):
    if s_to_t.setdefault(a, b) != b or t_to_s.setdefault(b, a) != a:
        return False
return True
```

`setdefault(a, b)` stores `b` the first time `a` is seen and returns whatever is
stored, so comparing the result to `b` checks consistency and records the pair
in one step.

With only `s_to_t`, `"badc"` and `"baba"` pass: `b -> b`, `a -> a`, `d -> b`,
`c -> a` never contradict each other. The reverse dict sees `b` already claimed
by `b` when `d` tries to take it.

## Useful variations

- A dict plus a `used` set of targets also works: a new key `a` whose target `b` is already in `used` is a collision.
- `len(set(s)) == len(set(t)) == len(set(zip(s, t)))` is a clever one-liner for the same check.
- Word Pattern is identical with `zip(pattern, words)` after `words = s.split()`.

## Interview relevance

Isomorphic Strings, Word Pattern and Find and Replace Pattern all reduce to
"is this a bijection?". Interviewers often hand you an input like `"badc"` /
`"baba"` to see whether you check both directions, so naming the reverse map up
front shows you understood the problem rather than pattern-matched it.
