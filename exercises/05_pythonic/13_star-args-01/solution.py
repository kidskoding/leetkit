def merge_settings(*layers, **overrides):
    merged = {}
    for layer in layers:
        merged.update(layer)
    merged.update(overrides)
    return merged
