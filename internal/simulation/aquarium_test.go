package simulation

import (
	"math"
	"math/rand"
	"testing"
)

func TestUpdateRetainsDeadFishFitnessForGeneration(t *testing.T) {
	deadGenome := validTestGenome()
	deadGenome.Color = newColor(10, 20, 30, OpaqueAlpha)
	survivorGenome := validTestGenome()
	aquarium := Simulation{config: DefaultConfig(), rng: rand.New(rand.NewSource(42)),
		fish: []Fish{
			{
				Genome:    deadGenome,
				Alive:     false,
				Age:       4,
				FoodEaten: 2,
				FishEaten: 1,
			},
			{
				Genome: survivorGenome,
				Health: DefaultConfig().MaxHealth,
				Alive:  true,
			},
		},
	}

	aquarium.Step(0)

	if got := len(aquarium.fish); got != 1 {
		t.Fatalf("living fish count = %d, want 1", got)
	}
	if got := aquarium.fish[0].Genome; got != survivorGenome {
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
	weakGenome := validTestGenome()
	weakGenome.Color = newColor(10, 20, 30, OpaqueAlpha)
	survivorGenome := validTestGenome()
	survivorGenome.Color = newColor(220, 230, 240, OpaqueAlpha)
	aquarium := Simulation{config: DefaultConfig(), rng: rand.New(rand.NewSource(42)),
		fish: []Fish{
			{
				Genome:    survivorGenome,
				Health:    DefaultConfig().MaxHealth,
				Alive:     true,
				Age:       100,
				FoodEaten: 3,
			},
		},
		generationElapsed: DefaultConfig().GenerationDuration,
		generationCandidates: []Candidate{
			{Genome: weakGenome, Fitness: 1},
		},
	}

	aquarium.Step(0)

	if got := len(aquarium.fish); got != DefaultConfig().PopulationSize {
		t.Fatalf("new population size = %d, want %d", got, DefaultConfig().PopulationSize)
	}
	if got := aquarium.fish[0].Genome; got != survivorGenome {
		t.Fatalf("elite genome = %+v, want surviving best genome %+v", got, survivorGenome)
	}
	if aquarium.generationCandidates != nil {
		t.Fatalf("generation candidates were not cleared: %+v", aquarium.generationCandidates)
	}

	for i, fish := range aquarium.fish {
		if !fish.Alive || fish.Health != DefaultConfig().MaxHealth || fish.Age != 0 || fish.FoodEaten != 0 || fish.FishEaten != 0 {
			t.Fatalf("fish %d runtime state was not reset: %+v", i, fish)
		}
	}
}

func TestGenerationRolloverPreservesElapsedOvershoot(t *testing.T) {
	aquarium := Simulation{config: DefaultConfig(), rng: rand.New(rand.NewSource(42)),
		fish: []Fish{
			{
				Genome: validTestGenome(),
				Health: DefaultConfig().MaxHealth,
				Alive:  true,
			},
		},
		generationElapsed: DefaultConfig().GenerationDuration - 0.25,
	}

	aquarium.Step(0.5)

	const want = float32(0.25)
	if math.Abs(float64(aquarium.generationElapsed-want)) > 0.0001 {
		t.Fatalf("generation elapsed = %v, want %v", aquarium.generationElapsed, want)
	}
}
