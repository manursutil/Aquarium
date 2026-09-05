package integration_test

import (
	"aquarium/internal/simulation"
	"reflect"
	"testing"
)

func TestTradeoffHistoryAndRollover(t *testing.T) {
	c := simulation.DefaultConfig()
	c.PredationSizeRatio = 1.2
	c.PopulationSize = 25
	first, err := simulation.New(c, 42)
	if err != nil {
		t.Fatal(err)
	}
	second, err := simulation.New(c, 42)
	if err != nil {
		t.Fatal(err)
	}
	completed := 0
	for step := 0; step < 20000 && completed < 10; step++ {
		first.Step(1.0 / 60)
		second.Step(1.0 / 60)
		h := first.History()
		if len(h) == completed {
			continue
		}
		completed++
		if len(h) != completed || h[completed-1].Generation != completed || h[completed-1].Evaluated != c.PopulationSize {
			t.Fatal("incorrect history accounting")
		}
		snap := first.Snapshot()
		if len(snap.Fish) != c.PopulationSize {
			t.Fatal("population not restored")
		}
		for _, f := range snap.Fish {
			g := f.Genome
			for _, trait := range []struct{ v, low, high float32 }{{g.Size, 10, 30}, {g.MaxSpeed, 2, 102}, {g.Vision, 50, 150}, {g.Metabolism, 0.05, 0.15}, {g.AttractionToFood, 0, 1}, {g.FearOfPredators, 0, 1}, {g.AttractionToOthers, 0, 1}} {
				if !(trait.v >= trait.low && trait.v <= trait.high) {
					t.Fatalf("illegal trait %+v", trait)
				}
			}
		}
		if !reflect.DeepEqual(h, second.History()) {
			t.Fatal("same-seed history diverged")
		}
	}
	if completed != 10 {
		t.Fatalf("completed only %d generations", completed)
	}
}
