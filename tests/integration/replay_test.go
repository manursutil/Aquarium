package integration_test

import (
	"aquarium/internal/simulation"
	"reflect"
	"testing"
)

func TestReplayAcrossManyGenerations(t *testing.T) {
	c := simulation.DefaultConfig()
	c.GenerationDuration = 1
	first, err := simulation.New(c, 42)
	if err != nil {
		t.Fatal(err)
	}
	second, err := simulation.New(c, 42)
	if err != nil {
		t.Fatal(err)
	}
	other, err := simulation.New(c, 43)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(first.Snapshot(), other.Snapshot()) {
		t.Fatal("different seeds produced identical worlds")
	}
	previous := first.Snapshot()
	for range 1200 {
		first.Step(0.05)
		second.Step(0.05)
		got := first.Snapshot()
		if !reflect.DeepEqual(got, second.Snapshot()) {
			t.Fatalf("replay diverged in generation %d", got.Generation)
		}
		if got.Generation != previous.Generation {
			if len(got.Fish) != c.PopulationSize {
				t.Fatal("rollover population changed")
			}
			for _, f := range got.Fish {
				if f.Health != c.MaxHealth || f.Age != 0 || f.Hunger != 0 || f.FoodEaten != 0 || f.FishEaten != 0 {
					t.Fatal("rollover did not reset runtime state")
				}
				if f.ID <= previous.Fish[0].ID {
					t.Fatal("fish identity was reused")
				}
			}
		}
		previous = got
	}
	if previous.Generation < 50 {
		t.Fatalf("only reached generation %d", previous.Generation)
	}
}

func TestSnapshotIsIndependent(t *testing.T) {
	s, err := simulation.New(simulation.DefaultConfig(), 42)
	if err != nil {
		t.Fatal(err)
	}
	original := s.Snapshot()
	changed := s.Snapshot()
	changed.Fish[0].Position.X = -999
	changed.Fish[0].Genome.Size = -999
	changed.Food[0].X = -999
	if !reflect.DeepEqual(original, s.Snapshot()) {
		t.Fatal("snapshot mutation changed simulation")
	}
	s.Step(0.1)
	if original.Fish[0].Age != 0 {
		t.Fatal("step mutated retained snapshot")
	}
}
