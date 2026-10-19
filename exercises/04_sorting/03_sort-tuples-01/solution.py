def rank_players(scores):
    return [name for name, _ in sorted(scores, key=lambda p: (-p[1], p[0]))]
