# Slicing

`seq[start:stop:step]` returns a new sequence with the selected elements. Any
part can be omitted, and negative indices count from the end.

```python
nums[-k:]   # last k elements
nums[:-k]   # everything except the last k
nums[::-1]  # reversed copy
```

Rotating right by `k` is just "last k elements, then the rest":

```python
k %= len(nums)
rotated = nums[-k:] + nums[:-k]
```

The modulo matters twice: it handles `k` larger than the length, and it turns
`k == len(nums)` into `k == 0`. With `k == 0`, `nums[-0:]` is the whole list
and `nums[:-0]` is empty, so the result is an unchanged copy.

Slices never raise IndexError; out-of-range bounds are clipped. They always
return a new list, so the input is left untouched.

## Useful variations

- `s == s[::-1]` checks whether a string is a palindrome.
- `nums[k:] + nums[:k]` rotates left by `k`.
- `nums[::2]` takes every other element; `nums[1::2]` takes the rest.
- `nums[:] = nums[-k:] + nums[:-k]` rotates in place by assigning to a slice.

## Interview relevance

Rotate Array, Rotate String, palindrome checks and "reverse the first k"
problems all reduce to slicing. Mention that a slice copies, so it costs O(k)
time and memory; if the interviewer asks for O(1) extra space, switch to the
three-reversal trick on the original list.
