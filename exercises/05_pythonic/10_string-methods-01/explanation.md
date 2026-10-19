# String methods: isalnum(), lower(), strip(), split()

Python strings carry the cleanup tools you need for messy input. Lowercase
everything and keep only letters and digits in one pass:

```python
cleaned = [ch.lower() for ch in s if ch.isalnum()]
return cleaned == cleaned[::-1]
```

Without them you end up hand-writing character ranges:

```python
if "a" <= ch <= "z" or "A" <= ch <= "Z" or "0" <= ch <= "9":
    ...
```

That is longer, easy to get wrong, and still needs a separate lowercasing step.

Every string method returns a new string (strings are immutable), so
`s.lower()` on its own does nothing; you have to use or assign the result.

## Useful variations

- `s.split()` with no argument splits on any run of whitespace and drops empty pieces, so `"  a   b ".split()` is `["a", "b"]`.
- `s.strip(".,!?")` trims those characters from both ends, handy for words with trailing punctuation.
- `ch.isalpha()`, `ch.isdigit()` and `ch.isspace()` test narrower character classes.
- Two pointers moving inward and skipping non-`isalnum()` characters solve the same problem in O(1) extra space.

## Interview relevance

Valid Palindrome (LeetCode 125) is a classic warm-up, and the same cleanup
appears in word counting, anagram checks and log parsing. Reaching straight for
`isalnum()` and `lower()` keeps the normalization obvious, and mentioning the
two-pointer follow-up shows you can trade memory for clarity on request.
