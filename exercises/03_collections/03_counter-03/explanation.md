# Counter subtraction

Subtracting one Counter from another subtracts counts key by key and then
drops every key whose result is zero or negative.

```python
from collections import Counter

Counter("aab") - Counter("ab")   # Counter({'a': 1})
Counter("ab") - Counter("aabb")  # Counter()
```

So `Counter(need) - Counter(have)` is exactly the shortfall. An empty Counter
is falsy, which makes the "does it fit?" check one line:

```python
not (Counter(note) - Counter(magazine))
```

Without Counter you would build a dict for the magazine, then loop over the
note, decrementing and returning `False` when a count goes negative.

## Useful variations

- `Counter(note) <= Counter(magazine)` (Python 3.10+) asks the same question with a subset comparison.
- `a.subtract(b)` updates in place and keeps zero and negative counts.
- `a + b`, `a & b` (min of counts), and `a | b` (max of counts) are the other multiset operators.

## Interview relevance

"Can X be built from Y?" appears as Ransom Note, Find Words That Can Be Formed
by Characters, and inventory-style checks. Knowing that `-` clips at zero
saves you a manual loop and makes the intent obvious to the interviewer.
