package simulation

import (
	"math/rand"
)

const (
	FoodEatenWeight float32 = 5
	FishEatenWeight float32 = 10
)

type Candidate struct {
	Genome    Genome
	Fitness   float32
	Age       float32
	FoodEaten int
	FishEaten int
}

func fitness(f Fish) float32 {
	return f.Age + float32(f.FoodEaten)*FoodEatenWeight + float32(f.FishEaten)*FishEatenWeight
}

func candidateFromFish(f Fish) Candidate {
	return Candidate{
		Genome:    f.Genome,
		Fitness:   fitness(f),
		Age:       f.Age,
		FoodEaten: f.FoodEaten,
		FishEaten: f.FishEaten,
	}
}

func makeCandidates(fish []Fish) []Candidate {
	candidates := make([]Candidate, len(fish))

	for i, f := range fish {
		candidates[i] = candidateFromFish(f)
	}

	return candidates
}

func selectParent(rng *rand.Rand, candidates []Candidate) Genome {
	best := candidates[rng.Intn(len(candidates))]

	for range 3 {
		candidate := candidates[rng.Intn(len(candidates))]
		if candidate.Fitness > best.Fitness {
			best = candidate
		}
	}

	return best.Genome
}

func avgChannel(a uint8, b uint8) uint8 {
	return uint8((int(a) + int(b)) / 2)
}

func mixTwoColors(a Color, b Color) Color {
	return Color{
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

func mutateFloat(rng *rand.Rand, config Config, value *float32, minV float32, maxV float32) {
	if rng.Float32() >= config.MutationRate {
		return
	}

	rangeSize := maxV - minV
	delta := (rng.Float32()*2 - 1) * rangeSize * config.MutationSize
	*value = max(minV, min(maxV, *value+delta))
}

func mutate(rng *rand.Rand, config Config, g *Genome) {
	if g == nil {
		return
	}

	mutateFloat(rng, config, &g.Size, MinFishSize, MinFishSize+FishSizeRange)
	mutateFloat(rng, config, &g.MaxSpeed, MinFishSpeed, MinFishSpeed+FishSpeedRange)
	mutateFloat(rng, config, &g.Vision, MinFishVision, MinFishVision+FishVisionRange)
	mutateFloat(rng, config, &g.Metabolism, MinMetabolism, MinMetabolism+MetabolismRange)

	mutateFloat(rng, config, &g.AttractionToFood, 0, 1)
	mutateFloat(rng, config, &g.FearOfPredators, 0, 1)
	mutateFloat(rng, config, &g.AttractionToOthers, 0, 1)

	if rng.Float32() < config.MutationRate {
		switch rng.Intn(3) {
		case 0:
			g.Color.R = jitterColorChannel(rng, g.Color.R)
		case 1:
			g.Color.G = jitterColorChannel(rng, g.Color.G)
		case 2:
			g.Color.B = jitterColorChannel(rng, g.Color.B)
		}
	}
}

func evolve(rng *rand.Rand, candidates []Candidate, config Config) []Genome {
	populationSize := config.PopulationSize
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

	eliteCount := min(config.EliteCount, populationSize)
	for i := range eliteCount {
		nextGeneration[i] = best.Genome
	}

	for i := eliteCount; i < populationSize; i++ {
		parent1 := selectParent(rng, candidates)
		parent2 := selectParent(rng, candidates)

		childGenome := crossover(parent1, parent2)
		mutate(rng, config, &childGenome)
		nextGeneration[i] = childGenome
	}

	return nextGeneration
}
