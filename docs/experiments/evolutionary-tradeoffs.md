# Evolutionary trade-offs

## Findings

**Scope:** 7 treatments × 10 paired seeds × 30 generations = 2,100 generations.
Survivors at the end of all 2,100 generations; zero extinctions.

| Comparison with zero-cost baseline | Mean difference | Seeds with that direction |
| --- | ---: | ---: |
| Tradeoff: final mean size | −1.27 | 8/10 smaller |
| Tradeoff: final mean speed | −5.57 | 7/10 slower |
| Tradeoff: mean survivors | −13.22 fish | 10/10 fewer |
| Tradeoff: mean best fitness | −47.51 | 10/10 lower |
| Movement cost: final mean speed | +10.83 | 9/10 faster |
| Body cost: final mean size | +0.83 | 8/10 larger |
| Sensing cost: final mean vision | +6.37 | 7/10 higher |

- **Combined rules:** smaller, slower populations, with lower survival and fitness.
- **Costs alone:** higher feeding totals and frequent increases in the charged trait.
  Mechanism to consider: stronger food steering at higher hunger.
- **Ratio 1.2 alone:** mean predation total 126.9 → 47.6 per run;
  more survivors in 10/10 pairs.
- **Strategy claim:** evidence for population shifts; insufficient evidence for
  distinct hunter and scavenger niches. Final mean size remains high at 26.97
  within the legal 10–30 range.

## Reproduce

Run from the repository root:

```sh
python3 scripts/tradeoff_experiment.py

# Recompute summaries from saved CSVs:
python3 scripts/tradeoff_experiment.py --summarize-only
```

| Resource | Contents |
| --- | --- |
| [Runner](../../scripts/tradeoff_experiment.py) | Treatment definitions and run commands |
| [Per-seed summary](../../results/tradeoffs/summary.csv) | Fitness, survival, meal totals, final trait ranges, first boundary generations |
| [Baseline, seed 42](../../results/tradeoffs/baseline-42.csv) | Generation-level control data |
| [Tradeoff, seed 42](../../results/tradeoffs/tradeoff-42.csv) | Generation-level combined treatment data |

Raw file pattern: `results/tradeoffs/<treatment>-<seed>.csv`.

## Experiment setup

| Parameter | Value |
| --- | --- |
| Seeds | 42–51, the same set for each treatment |
| Population | 50 fish per generation |
| Duration | 30 completed generations; 30 seconds per generation |
| Fixed step | 0.0166667 seconds |
| World | 800 × 600 |
| Food | 35 items; respawn on consumption |
| Health | Maximum 10; food healing 3; starvation damage 1/second |
| Evolution | Mutation rate 0.15; mutation size 0.10; elite count 1 |
| Fitness | `age_seconds + 5*food_eaten + 10*fish_eaten` |
| Evaluation population | Living and dead fish from the generation |

### Treatments

Energy weights: hunger/second. Acceleration: velocity change/second per steering response.

| Treatment | Movement | Size | Vision | Predation ratio | Acceleration |
| --- | ---: | ---: | ---: | ---: | ---: |
| Baseline | 0 | 0 | 0 | 1.0 | 0 |
| Movement | 0.20 | 0 | 0 | 1.0 | 0 |
| Body | 0 | 0.10 | 0 | 1.0 | 0 |
| Sensing | 0 | 0 | 0.10 | 1.0 | 0 |
| Costs | 0.20 | 0.10 | 0.10 | 1.0 | 0 |
| Ratio | 0 | 0 | 0 | 1.2 | 0 |
| Tradeoff | 0.20 | 0.10 | 0.10 | 1.2 | 30 |

## Results

| Column | Aggregation |
| --- | --- |
| Mean best / median best | Mean across seeds of each run's 30-generation mean / median best fitness |
| Survivors | Mean across seeds and generations |
| Food / predation | Mean whole-run meal total per seed |
| Size / speed / vision | Mean across seeds of generation-30 population means, including dead fish |

| Treatment | Mean best | Median best | Survivors | Food | Predation | Size | Speed | Vision |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| baseline | 175.54 | 173.25 | 45.75 | 25572.60 | 126.90 | 28.24 | 65.29 | 100.66 |
| movement | 195.32 | 193.00 | 45.77 | 30542.70 | 125.90 | 28.74 | 76.12 | 105.34 |
| body | 206.25 | 204.50 | 45.39 | 32768.30 | 136.50 | 29.07 | 65.72 | 101.80 |
| sensing | 196.89 | 195.50 | 45.56 | 30929.10 | 131.70 | 28.63 | 71.64 | 107.03 |
| costs | 246.90 | 246.25 | 45.74 | 43755.60 | 125.40 | 28.91 | 78.26 | 109.66 |
| ratio | 172.82 | 175.25 | 47.72 | 26606.00 | 47.60 | 27.80 | 68.81 | 99.01 |
| tradeoff | 128.03 | 126.75 | 32.52 | 13151.20 | 48.40 | 26.97 | 59.72 | 99.67 |


### Final individual ranges

Pooled minima and maxima across the ten seeds at generation 30:

| Treatment | Size | Speed | Vision |
| --- | --- | --- | --- |
| Baseline | 24.50–30.00 | 42.56–88.50 | 80.03–127.67 |
| Tradeoff | 22.78–30.00 | 39.37–84.62 | 82.62–122.97 |

### First boundary generation

Boundary = either legal trait limit. `None` = no boundary through generation 30.
CSV precision: four decimal places; values rounding to a limit count as a hit.

| Seed | Baseline size | Tradeoff size | Baseline vision | Tradeoff vision |
| ---: | --- | --- | --- | --- |
| 42 | None | None | None | None |
| 43 | 10 | 2 | 2 | None |
| 44 | None | 2 | None | None |
| 45 | 5 | 2 | None | None |
| 46 | 2 | 5 | None | None |
| 47 | 2 | 7 | None | None |
| 48 | 5 | None | None | None |
| 49 | 3 | 2 | None | None |
| 50 | 5 | 2 | None | 3 |
| 51 | None | 3 | None | None |

Speed: no boundary hits in either treatment. See the
[summary](../../results/tradeoffs/summary.csv) for the other five treatments.

## Rules

### Energy

For actual speed `v` and weights `Wspeed`, `Wsize`, `Wvision`:

```text
movement = Wspeed  × clamp(v / 102, 0, 1)²
size     = Wsize   × clamp((Size − 10) / 20, 0, 1)
vision   = Wvision × clamp((Vision − 50) / 100, 0, 1)

hunger increment = (Metabolism + movement + size + vision) × dt
```

| Constraint | Rule |
| --- | --- |
| Valid costs | Finite, nonnegative rates; representable combined hunger rate |
| Zero-cost control | Zero weights |
| Trait endpoints | Minimum size/vision: zero cost; maximum: configured weight |
| Movement | Actual speed; zero cost at rest |
| Time scaling | One multiplication by `dt` in the hunger phase |
| Genome | No changes from charging energy |
| Replay | Same seed, configuration, and step sequence |

Update order:

```text
Steering → movement → hunger/starvation → food → predation
         → dead-fish recording/removal → generation rollover
```

Food can reduce hunger after the current step's cost.

### Steering and predation

| Rule | Implementation |
| --- | --- |
| Steering limit | `BaseAcceleration × 10 / Size × dt` per response |
| Size example at acceleration 30 | Size 10: 30 velocity units/second; size 30: 10 |
| Limit scope | Each food, social, and predator response; existing target selection and behavior weights |
| Wandering | Existing rotation unchanged |
| Speed ceiling | `Genome.MaxSpeed` |
| Acceleration 0 | No acceleration cap |
| Predation size | `predator.Size > prey.Size` and `predator.Size >= prey.Size × ratio` |
| Predation color | Existing color-distance threshold, in both collision directions |
| Failed eligibility | No health or meal-counter change |
| Kill accounting | One meal per kill; no attacks by a dead fish later in the loop |
| Victim health | Lower bound zero |

Defaults: weights **0.20 / 0.10 / 0.10**, acceleration **30**, ratio **1.0**.
The combined experiment uses ratio **1.2**.

## Baseline and calibration archive

### Pre-change trace

Revision: `51f3b05`. Run this command at that revision to reproduce the original trace:

```sh
go run ./cmd/experiment -seed 42 -generations 100 -population 500 \
  -step 0.0166667 -output results/tradeoffs/pre-change-seed-42.csv
```

Configuration: 500 fish, 100 generations, 30 seconds/generation, 35 food items,
800 × 600 world, weights 0.20/0.10/0.10, ratio 1, original steering.
Fitness and step match the paired experiment.

| Metric | Generation 1 | Generation 100 |
| --- | ---: | ---: |
| Mean size | 19.90 | 29.27 |
| Mean speed | 52.45 | 49.85 |
| Mean vision | 99.77 | 93.18 |
| Best fitness | 209.82 | 4640.00 |
| Mean fitness | 14.30 | 2132.62 |
| Predation meals | 495 | 0 |

[Raw trace](../../results/tradeoffs/pre-change-seed-42.csv): 100 generations
with survivors. Size approaches its upper bound; speed and vision stay below
their upper bounds. Color groups: absent from the CSV.

Keep this 500-fish trace separate from the 50-fish paired comparison.

### Cost calibration before acceleration limits

Configuration: seed 42, 50 fish, 20 generations; other settings match the
pre-change trace. Generation-20 survivors: 50 in each treatment.

| Treatment and CSV | Final best fitness |
| --- | ---: |
| [Baseline](../../results/tradeoffs/calibration-baseline-42.csv) | 160 |
| [Movement](../../results/tradeoffs/calibration-movement-42.csv) | 180 |
| [Body](../../results/tradeoffs/calibration-body-42.csv) | 200 |
| [Sensing](../../results/tradeoffs/calibration-sensing-42.csv) | 190 |
| [Costs](../../results/tradeoffs/calibration-costs-42.csv) | 210 |

Calibration decision: retain the proposed weights for the paired runs.
The archive files use the original CSV schema, without trait extrema.

## Limits on interpretation

| Limit | Consequence |
| --- | --- |
| Ten seeds, 30 generations | Limited evidence for long-run strategy stability |
| Food respawn without a resource budget; fitness rewards meals | Advantage for active feeding |
| Stronger food steering at higher hunger | Energy costs can increase feeding and fitness |
| Combined changes to costs, ratio, and acceleration | No attribution of the combined effect to one coefficient |
| Social search includes the fish itself | Restricted interpretation of social behavior |
| “Predator avoidance” steers toward the selected other-color fish | Restricted interpretation of evasive strategies |

## Implementation and verification

Starting state: steps 2–5 had implementations at `51f3b05`.

| Checkpoint | Work |
| --- | --- |
| Baseline | Pre-change trace and cost calibration |
| Energy | Tests for configuration, normalization, breakdown, time scaling, and bounds |
| Predation | Shared size eligibility, collision guards, kill accounting |
| Steering | Size-dependent limits across three response paths |
| Integration | Ten nonzero-cost generations: history, population replacement, legal traits, replay |
| Reporting | Six CSV extrema columns; matching fixture; paired runs and summary |

Checks from the implementation run:

```sh
go test ./...
CGO_ENABLED=0 go test ./internal/simulation ./internal/report ./tests/integration
go vet ./...
```

Status: passed. Additional unit coverage: zero weights, invalid configuration,
stationary/maximal costs, speed squared, hunger scaling below saturation,
food/health bounds, color guards, and dead-fish collision handling.
