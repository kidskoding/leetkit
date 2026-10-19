from collections import namedtuple

Player = namedtuple("Player", ["name", "score"])


def parse_players(lines):
    players = []
    for line in lines:
        name, score = line.split(",")
        players.append(Player(name, int(score)))
    return players


def best_player(lines):
    players = parse_players(lines)
    if not players:
        return None
    return max(players, key=lambda player: player.score).name
