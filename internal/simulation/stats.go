package simulation

import "slices"

type TraitStats struct {
	Mean float32
	Min  float32
	Max  float32
}

type GenerationStats struct {
	BestFishID FishID
	Generation int
	Evaluated  int
	Survivors  int

	BestFitness   float32
	MeanFitness   float32
	MedianFitness float32

	FoodEaten int
	FishEaten int

	Size             TraitStats
	MaxSpeed         TraitStats
	Vision           TraitStats
	Metabolism       TraitStats
	FoodAttraction   TraitStats
	PredatorFear     TraitStats
	SocialAttraction TraitStats
}

func summarize(values []float32) TraitStats {
	if len(values) == 0 {
		return TraitStats{}
	}

	stats := TraitStats{
		Min: values[0],
		Max: values[0],
	}

	var total float32 = 0
	for _, value := range values {
		total += value

		if value < stats.Min {
			stats.Min = value
		}

		if value > stats.Max {
			stats.Max = value
		}
	}

	stats.Mean = total / float32(len(values))

	return stats
}

func median(values []float32) float32 {
	if len(values) == 0 {
		return 0
	}

	sorted := append([]float32(nil), values...)
	slices.Sort(sorted)

	middle := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[middle]
	}

	return (sorted[middle] + sorted[middle-1]) / 2
}

func summarizeGeneration(number int, survivors int, candidates []Candidate) GenerationStats {
	stats := GenerationStats{
		Generation: number,
		Evaluated:  len(candidates),
		Survivors:  survivors,
	}

	if len(candidates) == 0 {
		return stats
	}

	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.Fitness > best.Fitness {
			best = c
		}
	}
	stats.BestFishID = best.FishID

	fitnessValues := make([]float32, len(candidates))
	sizes := make([]float32, len(candidates))
	speeds := make([]float32, len(candidates))
	visions := make([]float32, len(candidates))
	metabolism := make([]float32, len(candidates))
	foodAttraction := make([]float32, len(candidates))
	predatorFear := make([]float32, len(candidates))
	socialAttraction := make([]float32, len(candidates))

	for i, candidate := range candidates {
		fitnessValues[i] = candidate.Fitness
		genome := candidate.Genome
		sizes[i] = genome.Size
		speeds[i] = genome.MaxSpeed
		visions[i] = genome.Vision
		metabolism[i] = genome.Metabolism
		foodAttraction[i] = genome.AttractionToFood
		predatorFear[i] = genome.FearOfPredators
		socialAttraction[i] = genome.AttractionToOthers
		stats.FoodEaten += candidate.FoodEaten
		stats.FishEaten += candidate.FishEaten
	}

	fitnessSummary := summarize(fitnessValues)
	stats.BestFitness = fitnessSummary.Max
	stats.MeanFitness = fitnessSummary.Mean
	stats.MedianFitness = median(fitnessValues)
	stats.Size = summarize(sizes)
	stats.MaxSpeed = summarize(speeds)
	stats.Vision = summarize(visions)
	stats.Metabolism = summarize(metabolism)
	stats.FoodAttraction = summarize(foodAttraction)
	stats.PredatorFear = summarize(predatorFear)
	stats.SocialAttraction = summarize(socialAttraction)

	return stats
}
