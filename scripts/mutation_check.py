"""Mutation-test exercise suites: inject small bugs into each reference
solution and report the ones its tests fail to catch.

Usage:
    python scripts/mutation_check.py [EXERCISE_DIR ...]

With no arguments every exercise is checked. Run it with a Python that has
pytest (for example .venv/bin/python). A surviving mutant is either a missing
test or an equivalent change that cannot alter the result; only the first kind
needs a new test.
"""

import argparse
import ast
import os
import shutil
import subprocess
import sys
import tempfile
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path

SWAPS = {
    ast.Lt: ast.LtE, ast.LtE: ast.Lt, ast.Gt: ast.GtE, ast.GtE: ast.Gt,
    ast.Eq: ast.NotEq, ast.NotEq: ast.Eq, ast.In: ast.NotIn, ast.NotIn: ast.In,
    ast.Add: ast.Sub, ast.Sub: ast.Add, ast.Mult: ast.Add, ast.FloorDiv: ast.Div,
    ast.And: ast.Or, ast.Or: ast.And,
}
TIMEOUT = 20


def mutation_sites(tree):
    for node in ast.walk(tree):
        if isinstance(node, ast.Compare):
            for i, op in enumerate(node.ops):
                if type(op) in SWAPS:
                    yield node, i
        elif isinstance(node, (ast.BinOp, ast.AugAssign, ast.BoolOp)) and type(node.op) in SWAPS:
            yield node, None
        elif isinstance(node, ast.Constant) and type(node.value) in (int, bool):
            yield node, None


def mutate(node, index):
    if isinstance(node, ast.Compare):
        node.ops[index] = SWAPS[type(node.ops[index])]()
        return "comparison"
    if isinstance(node, ast.Constant):
        old = node.value
        node.value = (not old) if isinstance(old, bool) else old + 1
        return f"{old!r} -> {node.value!r}"
    node.op = SWAPS[type(node.op)]()
    return "operator"


def mutants(source):
    count = sum(1 for _ in mutation_sites(ast.parse(source)))
    for k in range(count):
        tree = ast.parse(source)
        node, index = list(mutation_sites(tree))[k]
        what = mutate(node, index)
        yield f"line {node.lineno}: {what}", ast.unparse(tree)


def tests_pass(exercise, code):
    with tempfile.TemporaryDirectory() as tmp:
        Path(tmp, "starter.py").write_text(code)
        shutil.copy(exercise / "test_exercise.py", tmp)
        Path(tmp, "pytest.ini").write_text("[pytest]\n")
        try:
            result = subprocess.run(
                [sys.executable, "-B", "-m", "pytest", "-q", "-x", "-p", "no:cacheprovider", "test_exercise.py"],
                cwd=tmp, capture_output=True, timeout=TIMEOUT,
                env={**os.environ, "PYTEST_DISABLE_PLUGIN_AUTOLOAD": "1"},
            )
        except subprocess.TimeoutExpired:
            return False
        return result.returncode == 0


def check(exercise):
    source = (exercise / "solution.py").read_text()
    baseline = ast.unparse(ast.parse(source))
    total, survivors = 0, []
    for label, code in mutants(source):
        if code == baseline:
            continue
        total += 1
        if tests_pass(exercise, code):
            survivors.append(label)
    return exercise, total, survivors


def main():
    parser = argparse.ArgumentParser(description=__doc__.split("\n\n")[0])
    parser.add_argument("exercises", nargs="*", type=Path)
    args = parser.parse_args()
    root = Path(__file__).resolve().parent.parent / "exercises"
    dirs = args.exercises or sorted(p.parent for p in root.glob("*/*/solution.py"))

    with ThreadPoolExecutor(max_workers=os.cpu_count()) as pool:
        results = list(pool.map(check, dirs))

    total = sum(t for _, t, _ in results)
    survived = sum(len(s) for _, _, s in results)
    for exercise, _, survivors in results:
        for label in survivors:
            print(f"survived  {exercise.resolve().relative_to(root.parent)}  {label}")
    print(f"\n{len(results)} exercises, {total} mutants, {survived} survived")


if __name__ == "__main__":
    main()
