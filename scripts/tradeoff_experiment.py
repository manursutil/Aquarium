#!/usr/bin/env python3
"""Reproduce paired experiments (Python standard library + Go)."""

import csv
import statistics as stats
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "results/tradeoffs"
SEEDS = range(42, 52)
RUNS = {
    "baseline": (0, 0, 0, 1, 0),
    "movement": (0.2, 0, 0, 1, 0),
    "body": (0, 0.1, 0, 1, 0),
    "sensing": (0, 0, 0.1, 1, 0),
    "costs": (0.2, 0.1, 0.1, 1, 0),
    "ratio": (0, 0, 0, 1.2, 0),
    "tradeoff": (0.2, 0.1, 0.1, 1.2, 30),
}


def summarize():
    summaries = []
    for name in RUNS:
        for seed in SEEDS:
            path = OUT / f"{name}-{seed}.csv"
            with path.open() as source:
                rows = [
                    {k: float(v) for k, v in r.items()} for r in csv.DictReader(source)
                ]
            last = rows[-1]
            row = {
                "run": name,
                "seed": seed,
                "mean_best": stats.mean(r["best_fitness"] for r in rows),
                "median_best": stats.median(r["best_fitness"] for r in rows),
                "mean_survivors": stats.mean(r["survivors"] for r in rows),
                "food_total": sum(r["food_eaten"] for r in rows),
                "predation_total": sum(r["fish_eaten"] for r in rows),
                "nonextinct_generations": sum(r["survivors"] > 0 for r in rows),
            }
            for trait, low, high in [
                ("size", 10, 30),
                ("speed", 2, 102),
                ("vision", 50, 150),
            ]:
                row[f"final_{trait}_mean"] = last[f"{trait}_mean"]
                row[f"final_{trait}_min"] = last[f"{trait}_min"]
                row[f"final_{trait}_max"] = last[f"{trait}_max"]
                row[f"first_{trait}_boundary"] = next(
                    (
                        int(r["generation"])
                        for r in rows
                        if r[f"{trait}_min"] <= low or r[f"{trait}_max"] >= high
                    ),
                    "none",
                )
            summaries.append(row)
    with (OUT / "summary.csv").open("w", newline="") as dest:
        writer = csv.DictWriter(dest, fieldnames=summaries[0].keys())
        writer.writeheader()
        writer.writerows(summaries)
    for name in RUNS:
        rows = [r for r in summaries if r["run"] == name]
        print(
            name,
            {
                k: round(stats.mean(r[k] for r in rows), 3)
                for k in [
                    "mean_best",
                    "median_best",
                    "mean_survivors",
                    "food_total",
                    "predation_total",
                    "nonextinct_generations",
                    "final_size_mean",
                    "final_speed_mean",
                    "final_vision_mean",
                ]
            },
            flush=True,
        )


def main():
    import argparse

    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--summarize-only", action="store_true")
    args = parser.parse_args()

    OUT.mkdir(parents=True, exist_ok=True)

    if not args.summarize_only:
        import tempfile

        with tempfile.TemporaryDirectory(prefix="aquarium-tradeoff-") as tmp:
            binary = str(Path(tmp) / "experiment")
            subprocess.run(
                ["go", "build", "-o", binary, "./cmd/experiment"], cwd=ROOT, check=True
            )
            for name, (speed, size, vision, ratio, acceleration) in RUNS.items():
                for seed in SEEDS:
                    subprocess.run(
                        [
                            binary,
                            "-seed",
                            str(seed),
                            "-population",
                            "50",
                            "-generations",
                            "30",
                            "-duration",
                            "30",
                            "-food",
                            "35",
                            "-step",
                            "0.0166667",
                            "-speed-energy",
                            str(speed),
                            "-size-energy",
                            str(size),
                            "-vision-energy",
                            str(vision),
                            "-predation-ratio",
                            str(ratio),
                            "-acceleration",
                            str(acceleration),
                            "-output",
                            str(OUT / f"{name}-{seed}.csv"),
                        ],
                        check=True,
                        stdout=subprocess.DEVNULL,
                        stderr=subprocess.DEVNULL,
                    )
                print(f"Completed {name}: seeds 42–51", flush=True)

    summarize()


if __name__ == "__main__":
    main()
