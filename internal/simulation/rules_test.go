package simulation

import (
	"math/rand"
	"testing"
)

func TestMovementAndStarvationUseConfig(t *testing.T) {
	c := DefaultConfig()
	c.WorldWidth = 20
	c.WorldHeight = 30
	c.StarvationDamage = 2
	f := Fish{Position: Vector2{21, -1}, Health: 5, Hunger: 1, Alive: true, Genome: validTestGenome()}
	f.wrapEdges(c)
	if f.Position != (Vector2{0, 30}) {
		t.Fatalf("wrap: %+v", f.Position)
	}
	f.handleHunger(0.5, c)
	if f.Health != 4 {
		t.Fatalf("health=%v", f.Health)
	}
	f.Velocity = Vector2{3, 4}
	f.NoiseX = 0
	before := f.Position
	f.move(0.1)
	if distance(before, f.Position) < 0.49 {
		t.Fatal("fish did not move")
	}
}

func TestFoodCollisionAndExtinction(t *testing.T) {
	c := DefaultConfig()
	c.MaxHealth = 20
	c.FoodHealing = 7
	s := Simulation{config: c, rng: rand.New(rand.NewSource(1)), generation: 1}
	f := Fish{Genome: validTestGenome(), Health: 18, Hunger: 0.2, Alive: true}
	s.fish = []Fish{f}
	s.food = []Food{{HealAmount: c.FoodHealing}}
	s.handleFishFoodCollisions()
	if s.fish[0].Health != 20 || s.fish[0].Hunger != 0 || s.fish[0].FoodEaten != 1 {
		t.Fatalf("feeding: %+v", s.fish[0])
	}
	s.fish[0].Alive = false
	s.Step(0)
	if s.generation != 2 || len(s.fish) != c.PopulationSize {
		t.Fatal("extinction did not advance generation")
	}
}

func TestPredationUsesConfiguredSizeRatio(t *testing.T) {
	for _, ratio := range []float32{1, 1.2} {
		c := DefaultConfig()
		c.PredationSizeRatio = ratio
		predator := Fish{Genome: validTestGenome(), Alive: true, Health: 5}
		predator.Genome.Size = 21
		predator.Genome.Color = Color{255, 255, 255, 255}
		prey := Fish{Genome: validTestGenome(), Alive: true, Health: 3, Position: Vector2{1, 0}}
		prey.Genome.Size = 20
		prey.Genome.Color = Color{0, 0, 0, 255}
		s := Simulation{config: c, fish: []Fish{predator, prey}}
		s.handleFishFishCollisions()
		if ratio == 1 && (s.fish[1].Alive || s.fish[0].FishEaten != 1) {
			t.Fatal("original larger-fish predation was lost")
		}
		if ratio == 1.2 && (!s.fish[1].Alive || s.fish[1].Health != 3) {
			t.Fatal("configured ratio was ignored")
		}
	}
}
