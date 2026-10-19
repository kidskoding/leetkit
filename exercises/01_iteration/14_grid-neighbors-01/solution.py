DIRECTIONS = [(-1, 0), (1, 0), (0, -1), (0, 1)]


def find_peaks(heights):
    rows = len(heights)
    cols = len(heights[0]) if heights else 0
    peaks = []
    for r in range(rows):
        for c in range(cols):
            if all(
                heights[r][c] > heights[r + dr][c + dc]
                for dr, dc in DIRECTIONS
                if 0 <= r + dr < rows and 0 <= c + dc < cols
            ):
                peaks.append((r, c))
    return peaks
