package simulation

import "testing"

func TestPredationThresholds(t *testing.T) {
	for _, tc := range []struct {
		size, ratio float32
		want        bool
	}{{20, 1, false}, {20.1, 1, true}, {20, 1.2, false}, {23.99, 1.2, false}, {24, 1.2, true}, {25, 1.2, true}} {
		predator := Fish{Genome: Genome{Size: tc.size}}
		prey := Fish{Genome: Genome{Size: 20}}
		if canEat(predator, prey, tc.ratio) != tc.want {
			t.Errorf("size %g ratio %g", tc.size, tc.ratio)
		}
	}
}

func TestPredationCollisionGuards(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		for _, similar := range []bool{false, true} {
			for _, size := range []float32{23.99, 24} {
				c := DefaultConfig()
				c.PredationSizeRatio = 1.2
				predator := Fish{Alive: true, Health: 5, Genome: Genome{Size: size, Color: Color{255, 255, 255, 255}}}
				prey := Fish{Alive: true, Health: 3, Position: Vector2{1, 0}, Genome: Genome{Size: 20}}
				if similar {
					prey.Genome.Color = predator.Genome.Color
				}
				s := Simulation{config: c, fish: []Fish{predator, prey}}
				pi, qi := 0, 1
				if reverse {
					s.fish[0], s.fish[1] = s.fish[1], s.fish[0]
					pi, qi = 1, 0
				}
				s.handleFishFishCollisions()
				if similar || size < 24 {
					if s.fish[pi].Health != 5 || s.fish[qi].Health != 3 || s.fish[pi].FishEaten != 0 {
						t.Fatal("failed predation changed health/counter")
					}
				} else {
					if s.fish[qi].Alive || s.fish[pi].FishEaten != 1 {
						t.Fatal("successful predation not recorded")
					}
					s.fish[pi].attackPrey(&s.fish[qi], c)
					if s.fish[pi].FishEaten != 1 {
						t.Fatal("meal counted twice")
					}
				}
			}
		}
	}
}

func TestApplySteering(t *testing.T) {
	c := DefaultConfig()
	f := Fish{Genome: Genome{Size: MinFishSize, MaxSpeed: maximumPossibleSpeed}}
	desired := Vector2{100, 0}
	f.applySteering(desired, 0, c)
	near(t, length(f.Velocity), 0)
	f.applySteering(Vector2{}, 1, c)
	near(t, length(f.Velocity), 0)
	large := f
	large.Genome.Size = 30
	f.applySteering(desired, 0.5, c)
	near(t, length(f.Velocity), 15)
	large.applySteering(desired, 0.5, c)
	near(t, length(large.Velocity), 5)
	c.BaseAcceleration = 0
	f.applySteering(desired, 0.5, c)
	near(t, f.Velocity.X, 100)
	f.applySteering(Vector2{1000, 0}, 1, c)
	near(t, length(f.Velocity), maximumPossibleSpeed)
}

func TestSteeringPathsUseAcceleration(t *testing.T) {
	c := DefaultConfig()
	for _, mode := range []string{"food", "social", "predator"} {
		f := Fish{Alive: true, Hunger: 1, Genome: Genome{Size: 10, Vision: 150, MaxSpeed: 100, AttractionToFood: 1, AttractionToOthers: 1, FearOfPredators: 1}}
		other := Fish{Position: Vector2{50, 0}, Genome: f.Genome}
		switch mode {
		case "food":
			f.steerTowardFood([]Food{{Position: other.Position}}, 0.1, c)
		case "social":
			f.attractionToSimilarFish([]Fish{other}, 0.1, c)
		case "predator":
			other.Genome.Color = Color{255, 255, 255, 255}
			f.steerAwayFromPredators([]Fish{other}, 0.1, c)
		}
		near(t, length(f.Velocity), 3)
		var recorded Vector2
		switch mode {
		case "food":
			recorded = f.steering.Food
		case "social":
			recorded = f.steering.Social
		case "predator":
			recorded = f.steering.Predator
		}
		if recorded != f.Velocity {
			t.Fatalf("%s recorded %v, want applied change %v", mode, recorded, f.Velocity)
		}
	}
}

func TestDeadPredationVictimCannotAttackLaterFish(t *testing.T) {
	c := DefaultConfig()
	victim := Fish{Alive: true, Health: 3, Genome: Genome{Size: 20}}
	killer := Fish{Alive: true, Health: 5, Position: Vector2{1, 0}, Genome: Genome{Size: 30, Color: Color{255, 255, 255, 255}}}
	// After separating the first pair, the victim overlaps this fish but the killer does not.
	third := Fish{Alive: true, Health: 3, Position: Vector2{-25, 1}, Genome: Genome{Size: 10, Color: killer.Genome.Color}}
	s := Simulation{config: c, fish: []Fish{victim, killer, third}}
	s.handleFishFishCollisions()
	if s.fish[0].Health != 0 || s.fish[0].Alive || s.fish[0].FishEaten != 0 || s.fish[2].Health != 3 {
		t.Fatalf("dead fish attacked: %+v", s.fish)
	}
}
