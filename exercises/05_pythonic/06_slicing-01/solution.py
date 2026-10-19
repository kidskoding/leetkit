def rotate_right(nums, k):
    if not nums:
        return []
    k %= len(nums)
    return nums[-k:] + nums[:-k]
