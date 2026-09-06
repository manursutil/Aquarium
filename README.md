# Aquarium Genetic Algorithm

An evolutionary aquarium simulation written in Go with Raylib. Fish inherit
physical and behavioral traits, seek food, interact through color and size,
and produce new generations through selection, crossover, mutation, and
elitism.

The project is under development. The current version includes generation
tracking, fitness evaluation, bounded genome mutation, an on-screen HUD, and a
fish inspector.

## Run

Requires Go 1.27.1 or later.

```sh
go run ./cmd/aquarium
```

Space pauses or resumes. Keys `1`, `2`, and `3` select 1x, 5x, and 20x.
`R` restarts the displayed seed; `N` starts a new seed. Restarts preserve
configuration, pause/speed, and overlay preferences, and clear selection,
timing backlog, and summaries. Catch-up is capped at 40 steps per frame;
excess time is discarded after a stalled frame, so effective speed can drop
under load.

The HUD shows the seed, controls, and current settings. Completed generations
show a two-second summary using wall time, including while paused. Press `S`
to hide summaries at 20x. The fitness chart remains visible below the summary.

Click a fish to pin its inspector; Escape unpins it. Press `V` to toggle the
pinned fish's vision circle and `F` to toggle its steering vectors. Both overlays
start disabled. Green is food, gold is social attraction, magenta is predator
response, and the thicker white line is their combined response before wandering.
Component lines are 60 pixels long and the combined line is 80 pixels: they show
direction, not magnitude. A zero response draws no line. The snapshot retains
the actual velocity changes after behavioral weighting and acceleration/speed
limits; overlays do not change simulation behavior.

## Test

```sh
go test ./...
```

## Simulation API

`internal/simulation` owns world rules, food, movement, evolution, and generation
timing without importing Raylib. Create a world with
`simulation.New(simulation.DefaultConfig(), 42)`, advance it with `Step(dt)`
(seconds), and read `Snapshot()`. Snapshots own their slices and do not expose
live state. The same seed, configuration, and sequence of time steps reproduce
the same run.

`cmd/aquarium` owns the window, drawing, HUD, and hover inspector. It currently
starts with seed 42 and fixed 1/60-second simulation steps. Default tuning preserves the original
strictly-larger-fish predation rule; configuration can require a larger ratio.

Run the core and many-generation replay tests without graphics dependencies:

```sh
CGO_ENABLED=0 go test ./internal/simulation ./tests/integration
go vet ./...
```

## Evolutionary trade-offs

Movement, size, and vision contribute configurable hunger costs. Movement uses
actual speed squared; larger fish also have slower steering responses.
Predation can require a configurable size advantage. See the
[milestone 4 experiment](docs/experiments/milestone-4.md) for equations,
limitations, and paired results from ten seeds: the combined rules produced
smaller, slower populations, with fewer survivors and lower best fitness.

Reproduce the experiments with `python3 scripts/tradeoff_experiment.py`, or run
one treatment directly:

```sh
go run ./cmd/experiment -seed 42 -generations 30 -population 50 \
  -step 0.0166667 -speed-energy 0.2 -size-energy 0.1 -vision-energy 0.1 \
  -predation-ratio 1.2 -acceleration 30 -output results/tradeoff-42.csv
```

Use zero energy weights and `-acceleration 0 -predation-ratio 1` for the control.

## Lineage and color groups

The inspector shows birth generation, parents, elite provenance, and the living
population's color group/share. Pin a fish and press `A` to inspect three ancestry
levels. Press `B` to inspect the latest completed generation's best fish, even
after it leaves the live population. Escape closes ancestry and unpins selection.
Shared ancestors appear once with references from both branches.

Founders have no parents. Crossover children retain both selected parent IDs
(which can be equal). Elite copies receive a new ID and retain the best
candidate's genome with one source parent. Fish IDs are unique within a run;
restarting can reuse them. Completed records distinguish death from replacement
at rollover and retain final fitness; active records expose current fitness.
`Lineage(id)` and `LineageRecords()` return copies. `Config.LineageGenerations`
defaults to ten completed generations (minimum one), plus the current cohort.
Older IDs remain visible as “outside retained history” in ancestry.

Color groups are a reporting approximation for the roadmap's “species.” They
use an independent RGB Euclidean threshold of 60. Fish are processed in ID order
and assigned to the nearest fixed representative within the threshold, with ties
favoring the lower group ID. New groups start at 1 in each sample. Display colors
are the mean member RGB, rounded down. Members can be farther apart than the
threshold, and group IDs do not track a species across generations. Grouping
consumes no simulation randomness and does not affect steering or predation.

Export ancestry and group summaries alongside the unchanged generation CSV:

```sh
go run ./cmd/experiment -seed 42 -generations 12 -duration 0.2 -population 12 \
  -run-id demo-42 -output /tmp/generations.csv \
  -lineage-output /tmp/lineage.csv -groups-output /tmp/groups.csv
```

The command streams each completed cohort before retention can evict it, then
exports the new active cohort at the stopping point with an empty end reason
and current fitness. Each individual appears once. Group rows label `completed`
(all evaluated fish, including deaths) or `living` (the final active sample),
and include the threshold, population, share, mean fitness, and mean RGB.
Use `-color-group-threshold` to change export grouping independently of the
simulation. Combine runs by `(seed, run_id, fish_id)`; supply a distinct run ID
for each run, or use the default UTC timestamp. For byte-identical replay checks,
reuse an explicit run ID. Output files must have distinct paths.

