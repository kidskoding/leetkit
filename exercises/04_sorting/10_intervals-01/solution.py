def can_attend_all(intervals):
    meetings = sorted(intervals, key=lambda iv: iv[0])
    return all(prev[1] <= cur[0] for prev, cur in zip(meetings, meetings[1:]))
