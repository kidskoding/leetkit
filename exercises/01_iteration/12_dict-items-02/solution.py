def top_sellers(sales, n):
    return sorted(sales.items(), key=lambda kv: (-kv[1], kv[0]))[:n]
