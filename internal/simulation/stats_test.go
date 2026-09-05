package simulation

import (
	"slices"
	"testing"
)

func TestMedianDoesNotChangeInput(t *testing.T) {
	values := []float32{9, 1, 5, 3}
	wantInput := append([]float32(nil), values...)

	if got := median(values); got != 4 {
		t.Fatalf("median = %v, want 4", got)
	}

	if !slices.Equal(values, wantInput) {
		t.Fatalf("median changed input: %v", values)
	}
}

func TestSummarizeGeneration(t *testing.T) {
	candidates := []Candidate{
		{Genome: Genome{Size: 10, MaxSpeed: 20}, Fitness: 4, FoodEaten: 1},
		{Genome: Genome{Size: 14, MaxSpeed: 30}, Fitness: 8, FishEaten: 1},
	}

	got := summarizeGeneration(7, 1, candidates)
	if got.Generation != 7 || got.Evaluated != 2 || got.Survivors != 1 {
		t.Fatalf("population counts = %+v", got)
	}
	if got.BestFitness != 8 || got.MeanFitness != 6 || got.MedianFitness != 6 {
		t.Fatalf("fitness summary = %+v", got)
	}
	if got.FoodEaten != 1 || got.FishEaten != 1 {
		t.Fatalf("event totals = %+v", got)
	}
}
