package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	FoodEatenWeight float32 = 5
	FishEatenWeight float32 = 10
	MutationRate    float32 = 0.15
	MutationSize    float32 = 0.10
	EliteCount              = 1
)

type Candidate struct {
	Genome  Genome
	Fitness float32
}

func fitness(f Fish) float32 {
	return f.Age + float32(f.FoodEaten)*FoodEatenWeight + float32(f.FishEaten)*FishEatenWeight
}

func makeCandidates(fish []Fish) []Candidate {
	candidates := make([]Candidate, len(fish))

	for i, f := range fish {
		candidates[i] = Candidate{
			Genome:  f.Genome,
			Fitness: fitness(f),
		}
	}

	return candidates
}

func selectParent(candidates []Candidate) Genome {
	best := candidates[rand.Intn(len(candidates))]

	for range 3 {
		candidate := candidates[rand.Intn(len(candidates))]
		if candidate.Fitness > best.Fitness {
			best = candidate
		}
	}

	return best.Genome
}

func avgChannel(a uint8, b uint8) uint8 {
	return uint8((int(a) + int(b)) / 2)
}

func mixTwoColors(a rl.Color, b rl.Color) rl.Color {
	return rl.Color{
		R: avgChannel(a.R, b.R),
		G: avgChannel(a.G, b.G),
		B: avgChannel(a.B, b.B),
		A: OpaqueAlpha,
	}
}

func crossover(a Genome, b Genome) Genome {
	childGenome := Genome{
		Size:               (a.Size + b.Size) / 2,
		MaxSpeed:           (a.MaxSpeed + b.MaxSpeed) / 2,
		Vision:             (a.Vision + b.Vision) / 2,
		Metabolism:         (a.Metabolism + b.Metabolism) / 2,
		AttractionToFood:   (a.AttractionToFood + b.AttractionToFood) / 2,
		FearOfPredators:    (a.FearOfPredators + b.FearOfPredators) / 2,
		AttractionToOthers: (a.AttractionToOthers + b.AttractionToOthers) / 2,
		Color:              mixTwoColors(a.Color, b.Color),
	}

	return childGenome
}

func mutateFloat(value *float32, minV float32, maxV float32) {
	if rand.Float32() >= MutationRate {
		return
	}

	rangeSize := maxV - minV
	delta := (rand.Float32()*2 - 1) * rangeSize * MutationSize
	*value = max(minV, min(maxV, *value+delta))
}

func mutate(g *Genome) {
	if g == nil {
		return
	}

	mutateFloat(&g.Size, MinFishSize, MinFishSize+FishSizeRange)
	mutateFloat(&g.MaxSpeed, MinFishSpeed, MinFishSpeed+FishSpeedRange)
	mutateFloat(&g.Vision, MinFishVision, MinFishVision+FishVisionRange)
	mutateFloat(&g.Metabolism, MinMetabolism, MinMetabolism+MetabolismRange)

	mutateFloat(&g.AttractionToFood, 0, 1)
	mutateFloat(&g.FearOfPredators, 0, 1)
	mutateFloat(&g.AttractionToOthers, 0, 1)

	if rand.Float32() < MutationRate {
		switch rand.Intn(3) {
		case 0:
			g.Color.R = jitterColorChannel(g.Color.R)
		case 1:
			g.Color.G = jitterColorChannel(g.Color.G)
		case 2:
			g.Color.B = jitterColorChannel(g.Color.B)
		}
	}
}

func evolve(candidates []Candidate, populationSize int) []Genome {
	if populationSize <= 0 || len(candidates) == 0 {
		return nil
	}

	nextGeneration := make([]Genome, populationSize)

	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if candidate.Fitness > best.Fitness {
			best = candidate
		}
	}

	eliteCount := min(EliteCount, populationSize)
	for i := range eliteCount {
		nextGeneration[i] = best.Genome
	}

	for i := eliteCount; i < populationSize; i++ {
		parent1 := selectParent(candidates)
		parent2 := selectParent(candidates)

		childGenome := crossover(parent1, parent2)
		mutate(&childGenome)
		nextGeneration[i] = childGenome
	}

	return nextGeneration
}
