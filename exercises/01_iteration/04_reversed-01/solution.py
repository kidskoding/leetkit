def last_index(events, user):
    for i in reversed(range(len(events))):
        if events[i] == user:
            return i
    return -1
