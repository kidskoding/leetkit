def orders_by_customer(orders):
    groups = {}
    for customer, order_id in orders:
        groups.setdefault(customer, []).append(order_id)
    return groups
