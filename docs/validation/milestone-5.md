# Milestone 5 validation

Implementation and automated checks completed on 2026-09-05. The manual
interaction gate remains pending because Computer Use rejected access to the
local test app: `Computer Use was not approved to use AquariumM5`.

## Roadmap steps

1. View state owns pause, speed, overlays, selection, and summary preferences.
2. Application timing uses fixed 1/60-second steps. The catch-up cap is 40,
   increased from the example's 10 so 20x is reachable at 60 FPS. A float64
   accumulator avoids rounding a 20-step budget down to 19. Pausing preserves
   the fractional remainder and performs no steps.
3. R recreates the current seed; N chooses a different application-owned seed.
   Both reset history, selection, summaries, and the accumulator while preserving
   simulation configuration and display preferences.
4. Pinned selection resolves by ID and clears when the ID disappears.
5. Copied steering data records applied velocity changes and resets each step.
6. Optional selected-fish overlays use distinct colors and fixed display lengths.
7. Both inspectors include ID, generation, actual speed, runtime values, and genome.
8. Latest-generation summaries last two seconds of rendered time. S optionally
   hides them at 20x. Multiple completions between frames select the latest record.
9. HUD includes the seed and speed/pause state. A measured-width bottom controls
   panel includes overlay, pinning, restart, and summary settings. The chart sits
   below the summary and above the controls at the default 800 by 600 size.
10. Automated checks pass; manual keyboard and visual checks remain pending.

## Passed checks

- `go fmt ./...`
- `CGO_ENABLED=0 go test ./internal/simulation ./internal/report ./tests/integration`
- `go test ./...`
- `go vet ./...`
- `git diff --check`
- Built executable and launched the Raylib window successfully at 800 by 600.

Tests cover fixed-step budgets at 1x/5x/20x, zero elapsed time, capped backlog,
pause, same-seed initial snapshots and replay history, different-seed restart,
selection after compaction/rollover/restart, summary activation/expiry, steering
recording, and snapshot/overlay isolation over multiple generations.

## Remaining manual gate

Follow the eight-action window sequence in `docs/dev/portfolio-roadmap.md`,
Milestone 5, Step 10. Also check 20x and the S summary toggle. Window launch alone
does not verify key handling, mouse behavior, or visual readability.
