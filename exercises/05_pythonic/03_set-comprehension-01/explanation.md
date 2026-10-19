# Set comprehensions

A set comprehension builds a set, so duplicates disappear automatically.

```python
domains = {email.rpartition("@")[2].lower() for email in emails}
```

It replaces the loop that creates an empty set and adds to it:

```python
domains = set()
for email in emails:
    domains.add(email.rpartition("@")[2].lower())
```

Normalize inside the expression (lowercase, strip, round) before the value
lands in the set; that is what makes "Gmail.com" and "gmail.com" collapse into
one entry.

Watch out: `{}` is an empty dict, not an empty set. Write `set()` for an empty
set.

## Useful variations

- `{ch for ch in s if ch.isalpha()}` collects distinct letters.
- `len({x % k for x in nums})` counts distinct remainders.
- `{tuple(sorted(word)) for word in words}` deduplicates anagram groups.
- `"a@b.com".rpartition("@")` returns `("a", "@", "b.com")`; index `[2]` is the
  part after the last separator.

## Interview relevance

"Count the unique ..." and "how many distinct ..." questions often reduce to
the size of a set built from normalized values. LeetCode's Unique Email
Addresses is exactly this pattern. A set comprehension keeps the
normalization rule in one visible place, which makes it easy to talk through.
