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

## Test

```sh
go test ./...
```
