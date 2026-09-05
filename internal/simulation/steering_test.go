package simulation

import "testing"

func TestSteeringSnapshotResetsAndCopiesAppliedResponse(t *testing.T) {
	c := DefaultConfig()
	c.PopulationSize = 1
	c.FoodCount = 0
	s, err := New(c, 42)
	if err != nil {
		t.Fatal(err)
	}
	f := &s.fish[0]
	f.Position = Vector2{100, 100}
	f.Velocity = Vector2{0, 5}
	f.Hunger = 1
	f.Genome.Vision = 150
	f.Genome.MaxSpeed = 100
	f.Genome.AttractionToFood = 1
	s.food = []Food{{Position: Vector2{200, 100}}}
	// Stale debug data must be cleared even for rules with no target.
	f.steering = SteeringSnapshot{Social: Vector2{1, 2}, Predator: Vector2{3, 4}}
	expected := *f
	expected.steerTowardFood(s.food, 0.1, c)
	want := subtract(expected.Velocity, f.Velocity)
	if want == (Vector2{}) {
		t.Fatal("fixture did not produce steering")
	}
	s.Step(0.1)
	got := s.Snapshot().Fish[0].Steering
	if got != (SteeringSnapshot{Food: want, Final: want}) {
		t.Fatalf("steering snapshot = %+v, want food and final %v", got, want)
	}
	s.food = nil
	s.Step(0.1)
	if got := s.Snapshot().Fish[0].Steering; got != (SteeringSnapshot{}) {
		t.Fatalf("target disappeared but steering persisted: %+v", got)
	}
}

func TestZeroBehaviorWeightRecordsNoResponse(t *testing.T) {
	f := Fish{Hunger: 1, Genome: Genome{Size: 10, Vision: 150, MaxSpeed: 100}}
	f.steerTowardFood([]Food{{Position: Vector2{50, 0}}}, 0.1, DefaultConfig())
	if f.steering.Food != (Vector2{}) {
		t.Fatalf("zero attraction recorded an applied response: %v", f.steering.Food)
	}
}
