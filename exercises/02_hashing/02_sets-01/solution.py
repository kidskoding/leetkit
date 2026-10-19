def has_duplicate(order_ids):
    seen = set()
    for order_id in order_ids:
        if order_id in seen:
            return True
        seen.add(order_id)
    return False
