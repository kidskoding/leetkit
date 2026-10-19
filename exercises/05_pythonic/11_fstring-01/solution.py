def receipt_line(item, qty, unit_price):
    total = qty * unit_price
    return f"{item:<10} {qty:>3} {total:>8.2f}"
