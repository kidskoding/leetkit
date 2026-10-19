def count_correct(answers, key):
    return sum(given == correct for given, correct in zip(answers, key))
