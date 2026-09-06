#!/usr/bin/env python3
"""Run Go benchmarks and save raw output, metadata, CSV, and a chart."""
import argparse
import csv
import os
from pathlib import Path
import platform
import re
import statistics
import subprocess

import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt

ROOT = Path(__file__).resolve().parents[1]
PATTERN = re.compile(r"^BenchmarkCohort60Steps/population_(\d+)(?:-\d+)?\s+\d+\s+([\d.]+) ns/op\s+([\d.]+) B/op\s+([\d.]+) allocs/op", re.M)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", type=Path, default=ROOT / "results/benchmarks")
    parser.add_argument("--count", type=int, default=5)
    args = parser.parse_args()
    if args.count < 1:
        parser.error("--count must be positive")
    command = ["go", "test", "./internal/simulation", "-run=^$",
               "-bench=^BenchmarkCohort60Steps$", "-benchmem", "-benchtime=1s",
               f"-count={args.count}", "-cpu=1"]
    result = subprocess.run(command, cwd=ROOT, env={**os.environ, "CGO_ENABLED": "0"},
                            check=True, capture_output=True, text=True)
    rows = [(int(n), float(ns), float(size), float(allocs))
            for n, ns, size, allocs in PATTERN.findall(result.stdout)]
    if any(sum(row[0] == n for row in rows) != args.count for n in (50, 100, 250, 500)):
        raise ValueError(f"Incomplete benchmark output:\n{result.stdout}")
    args.output.mkdir(parents=True, exist_ok=True)
    (args.output / "raw.txt").write_text(result.stdout)
    version = subprocess.check_output(["go", "version"], cwd=ROOT, text=True).strip()
    (args.output / "environment.txt").write_text(
        f"{platform.platform()}\n{version}\nCGO_ENABLED=0\n{' '.join(command)}\n")
    with (args.output / "samples.csv").open("w", newline="") as dest:
        writer = csv.writer(dest)
        writer.writerow(["initial_population", "ns_per_op", "bytes_per_op", "allocs_per_op"])
        writer.writerows(rows)
    fig, ax = plt.subplots(figsize=(7, 4), layout="constrained")
    for n in (50, 100, 250, 500):
        values = [ns / 1e6 for population, ns, _, _ in rows if population == n]
        ax.scatter([n] * len(values), values, color="#287d8e", alpha=.6)
        ax.plot([n - 12, n + 12], [statistics.median(values)] * 2, color="#ba4934")
    ax.set(xlabel="Initial fish population", ylabel="Milliseconds / workload",
           title="Initialization + 60 steps · seed 42 · one CPU\nDots: repeats; red bars: medians")
    ax.set_xticks([50, 100, 250, 500])
    ax.grid(axis="y", alpha=.2)
    fig.savefig(args.output / "benchmark.png", dpi=160)
    plt.close(fig)
    print(result.stdout)
    print(args.output / "benchmark.png")


if __name__ == "__main__":
    main()
