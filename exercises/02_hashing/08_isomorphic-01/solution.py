def is_isomorphic(s, t):
    if len(s) != len(t):
        return False
    s_to_t = {}
    t_to_s = {}
    for a, b in zip(s, t):
        if s_to_t.setdefault(a, b) != b or t_to_s.setdefault(b, a) != a:
            return False
    return True
