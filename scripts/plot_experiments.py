#!/usr/bin/env python3
"""Plot saved paired runs; no simulation rerun required."""
import argparse
import csv
from pathlib import Path

import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from tradeoff_experiment import ROOT, RUNS, SEEDS


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--input", type=Path, default=ROOT / "results/tradeoffs")
    parser.add_argument("--output", type=Path, default=ROOT / "docs/images/experiments.png")
    args = parser.parse_args()
    data = {}
    for treatment in RUNS:
        data[treatment] = []
        for seed in SEEDS:
            with (args.input / f"{treatment}-{seed}.csv").open() as source:
                rows = list(csv.DictReader(source))
            if len(rows) != 30 or [int(r["generation"]) for r in rows] != list(range(1, 31)):
                raise ValueError(f"Expected generations 1–30: {treatment}, seed {seed}")
            data[treatment].append(rows)
    fig, axes = plt.subplots(1, 2, figsize=(11, 4), layout="constrained")
    labels = list(RUNS)
    for ax, column, title in zip(axes, ["speed_mean", "survivors"],
                                  ["Generation 30: mean maximum speed", "Mean survivors over 30 generations"]):
        for x, label in enumerate(labels):
            values = [float(rows[-1][column]) if column == "speed_mean" else
                      sum(float(r[column]) for r in rows) / len(rows)
                      for rows in data[label]]
            ax.scatter([x + (i - 4.5) * .035 for i in range(10)], values,
                       s=22, alpha=.7, color="#287d8e")
            ax.plot([x - .25, x + .25], [sum(values) / len(values)] * 2,
                    color="#ba4934", linewidth=2)
        ax.set_xticks(range(len(labels)), labels, rotation=35, ha="right")
        ax.set_title(title)
        ax.set_ylabel("Maximum speed (world units/s)" if column == "speed_mean" else "Fish")
        ax.grid(axis="y", alpha=.2)
    fig.suptitle("10 paired seeds · 50 fish · dots: seeds; red bars: means")
    args.output.parent.mkdir(parents=True, exist_ok=True)
    fig.savefig(args.output, dpi=160)
    plt.close(fig)
    print(args.output)


if __name__ == "__main__":
    main()
