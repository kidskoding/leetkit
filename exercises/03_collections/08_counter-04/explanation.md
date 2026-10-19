# Sliding a Counter window

Instead of building a fresh Counter for every substring, keep one Counter for
the window and adjust two counts per step: the letter entering on the right
and the letter leaving on the left.

```python
from collections import Counter

window = Counter()
for i, ch in enumerate(s):
    window[ch] += 1
    if i >= k:
        window[s[i - k]] -= 1
    if window == target:
        starts.append(i - k + 1)
```

The clunky version rebuilds the whole window each time, which is O(k) work per
position:

```python
if Counter(s[i:i + k]) == target: ...
```

Since Python 3.10, Counter equality treats missing keys and zero counts the
same, so a letter that drops to 0 does not break `window == target`.

## Useful variations

- On older Pythons, delete keys that hit zero (`if window[c] == 0: del window[c]`) before comparing.
- Track a `matches` count of letters whose counts agree to make each comparison O(1).
- The same window answers Permutation in String: return `True` on the first match.

## Interview relevance

Fixed-size anagram windows appear as Find All Anagrams in a String,
Permutation in String, and substring-with-concatenation variants. Updating
counts incrementally is the step that turns a quadratic scan into a linear
one, and Counter keeps the bookkeeping short enough to write on a whiteboard.
