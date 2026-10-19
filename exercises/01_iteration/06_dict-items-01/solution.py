def over_budget(spending, limit):
    return [category for category, amount in spending.items() if amount > limit]
