package simulation

import (
	"math/rand"
	"testing"
)

func TestFitnessRewardsSurvivalAndFeeding(t *testing.T) {
	fish := Fish{
		Age:       12,
		FoodEaten: 2,
		FishEaten: 1,
	}

	const want = float32(32)
	if got := fitness(fish); got != want {
		t.Fatalf("fitness = %v, want %v", got, want)
	}
}

func TestCrossoverAveragesColorChannelsWithoutOverflow(t *testing.T) {
	a := validTestGenome()
	b := validTestGenome()
	a.Color = newColor(240, 70, 250, 10)
	b.Color = newColor(245, 210, 250, 20)

	got := crossover(a, b).Color
	want := newColor(242, 140, 250, OpaqueAlpha)
	if got != want {
		t.Fatalf("crossover color = %+v, want %+v", got, want)
	}
}

func TestMutationKeepsGenomeWithinValidRanges(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	genome := validTestGenome()
	changed := false

	for range 2_000 {
		before := genome
		mutate(rng, DefaultConfig(), &genome)
		assertGenomeWithinBounds(t, genome)
		if genome != before {
			changed = true
		}
	}

	if !changed {
		t.Fatal("genome never changed after 2,000 mutation attempts")
	}
}

func TestEvolveReturnsRequestedPopulationSize(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	candidates := []Candidate{{Genome: validTestGenome(), Fitness: 1}}

	const populationSize = 7
	if got := len(evolve(rng, candidates, testConfig(populationSize))); got != populationSize {
		t.Fatalf("next generation size = %d, want %d", got, populationSize)
	}
}

func TestEvolvePreservesBestGenomeAsElite(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	weakGenome := validTestGenome()
	bestGenome := validTestGenome()
	bestGenome.Color = newColor(1, 2, 3, OpaqueAlpha)
	candidates := []Candidate{
		{Genome: weakGenome, Fitness: 10},
		{Genome: bestGenome, Fitness: 100},
	}

	nextGeneration := evolve(rng, candidates, testConfig(5))

	if got := nextGeneration[0].Genome; got != bestGenome {
		t.Fatalf("elite genome = %+v, want %+v", got, bestGenome)
	}
}

func validTestGenome() Genome {
	return Genome{
		Size:               MinFishSize + FishSizeRange/2,
		MaxSpeed:           MinFishSpeed + FishSpeedRange/2,
		Vision:             MinFishVision + FishVisionRange/2,
		Metabolism:         MinMetabolism + MetabolismRange/2,
		AttractionToFood:   0.5,
		FearOfPredators:    0.5,
		AttractionToOthers: 0.5,
		Color:              newColor(100, 120, 140, OpaqueAlpha),
	}
}

func assertGenomeWithinBounds(t *testing.T, genome Genome) {
	t.Helper()

	checks := []struct {
		name     string
		value    float32
		min, max float32
	}{
		{name: "size", value: genome.Size, min: MinFishSize, max: MinFishSize + FishSizeRange},
		{name: "max speed", value: genome.MaxSpeed, min: MinFishSpeed, max: MinFishSpeed + FishSpeedRange},
		{name: "vision", value: genome.Vision, min: MinFishVision, max: MinFishVision + FishVisionRange},
		{name: "metabolism", value: genome.Metabolism, min: MinMetabolism, max: MinMetabolism + MetabolismRange},
		{name: "attraction to food", value: genome.AttractionToFood, min: 0, max: 1},
		{name: "fear of predators", value: genome.FearOfPredators, min: 0, max: 1},
		{name: "attraction to others", value: genome.AttractionToOthers, min: 0, max: 1},
	}

	for _, check := range checks {
		if check.value < check.min || check.value > check.max {
			t.Errorf("%s = %v, want value in [%v, %v]", check.name, check.value, check.min, check.max)
		}
	}

	if genome.Color.A != OpaqueAlpha {
		t.Errorf("color alpha = %d, want %d", genome.Color.A, OpaqueAlpha)
	}
}

func testConfig(n int) Config { c := DefaultConfig(); c.PopulationSize = n; return c }
