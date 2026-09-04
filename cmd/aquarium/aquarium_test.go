package main

import (
	"math"
	"math/rand"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestUpdateRetainsDeadFishFitnessForGeneration(t *testing.T) {
	deadGenome := validTestGenome()
	deadGenome.Color = rl.NewColor(10, 20, 30, OpaqueAlpha)
	survivorGenome := validTestGenome()
	aquarium := Aquarium{
		Fish: []Fish{
			{
				Genome:    deadGenome,
				Alive:     false,
				Age:       4,
				FoodEaten: 2,
				FishEaten: 1,
			},
			{
				Genome: survivorGenome,
				Health: MaxHealth,
				Alive:  true,
			},
		},
	}

	aquarium.update(0)

	if got := len(aquarium.Fish); got != 1 {
		t.Fatalf("living fish count = %d, want 1", got)
	}
	if got := aquarium.Fish[0].Genome; got != survivorGenome {
		t.Fatalf("remaining genome = %+v, want survivor %+v", got, survivorGenome)
	}
	if got := len(aquarium.generationCandidates); got != 1 {
		t.Fatalf("candidate count = %d, want 1", got)
	}

	wantCandidate := Candidate{Genome: deadGenome, Fitness: 24}
	if got := aquarium.generationCandidates[0]; got != wantCandidate {
		t.Fatalf("dead fish candidate = %+v, want %+v", got, wantCandidate)
	}
}

func TestGenerationRolloverIncludesSurvivorsAndResetsPopulation(t *testing.T) {
	rand.Seed(4)
	weakGenome := validTestGenome()
	weakGenome.Color = rl.NewColor(10, 20, 30, OpaqueAlpha)
	survivorGenome := validTestGenome()
	survivorGenome.Color = rl.NewColor(220, 230, 240, OpaqueAlpha)
	aquarium := Aquarium{
		Fish: []Fish{
			{
				Genome:    survivorGenome,
				Health:    MaxHealth,
				Alive:     true,
				Age:       100,
				FoodEaten: 3,
			},
		},
		generationElapsed: GenerationDuration,
		generationCandidates: []Candidate{
			{Genome: weakGenome, Fitness: 1},
		},
	}

	aquarium.update(0)

	if got := len(aquarium.Fish); got != InitialFishCount {
		t.Fatalf("new population size = %d, want %d", got, InitialFishCount)
	}
	if got := aquarium.Fish[0].Genome; got != survivorGenome {
		t.Fatalf("elite genome = %+v, want surviving best genome %+v", got, survivorGenome)
	}
	if aquarium.generationCandidates != nil {
		t.Fatalf("generation candidates were not cleared: %+v", aquarium.generationCandidates)
	}

	for i, fish := range aquarium.Fish {
		if !fish.Alive || fish.Health != MaxHealth || fish.Age != 0 || fish.FoodEaten != 0 || fish.FishEaten != 0 {
			t.Fatalf("fish %d runtime state was not reset: %+v", i, fish)
		}
	}
}

func TestGenerationRolloverPreservesElapsedOvershoot(t *testing.T) {
	rand.Seed(5)
	aquarium := Aquarium{
		Fish: []Fish{
			{
				Genome: validTestGenome(),
				Health: MaxHealth,
				Alive:  true,
			},
		},
		generationElapsed: GenerationDuration - 0.25,
	}

	aquarium.update(0.5)

	const want = float32(0.25)
	if math.Abs(float64(aquarium.generationElapsed-want)) > 0.0001 {
		t.Fatalf("generation elapsed = %v, want %v", aquarium.generationElapsed, want)
	}
}
