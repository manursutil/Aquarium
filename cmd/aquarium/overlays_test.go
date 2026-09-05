package main

import (
	"math"
	"reflect"
	"testing"

	"aquarium/internal/simulation"
)

func TestOverlaySelectionAndGeometry(t *testing.T) {
	fish := []simulation.FishSnapshot{
		{ID: 1},
		{ID: 7, Genome: simulation.Genome{Vision: 120}, Steering: simulation.SteeringSnapshot{
			Food:   simulation.Vector2{X: 3, Y: 4},
			Social: simulation.Vector2{X: -2},
			Final:  simulation.Vector2{X: 1, Y: 4},
		}},
	}
	for _, vision := range []bool{false, true} {
		for _, steering := range []bool{false, true} {
			view := ViewState{HasSelectedFish: true, SelectedFishID: 7, ShowVision: vision, ShowSteering: steering}
			overlay := buildFishOverlay(fish, view)
			if overlay.ShowVision != vision || (len(overlay.Vectors) > 0) != steering {
				t.Fatalf("flags not respected: %+v", overlay)
			}
			if vision || steering {
				if overlay.Fish.ID != 7 || overlay.Fish.Genome.Vision != 120 {
					t.Fatal("overlay did not follow pinned ID")
				}
			}
			if steering {
				if len(overlay.Vectors) != 4 {
					t.Fatal("missing steering components")
				}
				for i, vector := range overlay.Vectors {
					want := 60.0
					if i == 0 {
						want = 80
					}
					if i == 3 {
						want = 0
					}
					if math.Abs(math.Hypot(float64(vector.Vector.X), float64(vector.Vector.Y))-want) > 0.001 {
						t.Fatalf("vector %d has wrong display length: %v", i, vector.Vector)
					}
					for j := 0; j < i; j++ {
						if vector.Color == overlay.Vectors[j].Color {
							t.Fatal("component colors are not distinct")
						}
					}
				}
				if overlay.Vectors[0].Thickness <= overlay.Vectors[1].Thickness {
					t.Fatal("final vector is not heavier")
				}
				if overlay.Vectors[1].Vector.X <= 0 || overlay.Vectors[1].Vector.Y <= 0 || overlay.Vectors[2].Vector.X >= 0 {
					t.Fatal("normalization changed direction")
				}
			}
		}
	}
	for _, view := range []ViewState{
		{ShowVision: true, ShowSteering: true},
		{HasSelectedFish: true, SelectedFishID: 99, ShowVision: true, ShowSteering: true},
	} {
		if got := buildFishOverlay(fish, view); got.ShowVision || len(got.Vectors) != 0 {
			t.Fatal("overlay drawn without an existing pinned fish")
		}
	}
}

func TestOverlayReplayAndSnapshotIsolation(t *testing.T) {
	c := simulation.DefaultConfig()
	c.PopulationSize = 8
	c.FoodCount = 5
	c.GenerationDuration = 0.5
	for _, vision := range []bool{false, true} {
		for _, steering := range []bool{false, true} {
			control, err := simulation.New(c, 42)
			if err != nil {
				t.Fatal(err)
			}
			sim, err := simulation.New(c, 42)
			if err != nil {
				t.Fatal(err)
			}
			for frame := 0; frame < 180; frame++ {
				control.Step(fixedStep)
				sim.Step(fixedStep)
				snapshot := sim.Snapshot()
				view := ViewState{HasSelectedFish: true, SelectedFishID: snapshot.Fish[0].ID, ShowVision: vision, ShowSteering: steering}
				// Exercise the same preparation path used by drawAquarium, headlessly.
				overlay := buildFishOverlay(snapshot.Fish, view)
				overlay.Fish.Steering.Final.X = -999
				for i := range snapshot.Fish {
					snapshot.Fish[i].Steering = simulation.SteeringSnapshot{
						Food: simulation.Vector2{X: -999}, Social: simulation.Vector2{Y: -999},
						Predator: simulation.Vector2{X: 999}, Final: simulation.Vector2{Y: 999},
					}
				}
				if !reflect.DeepEqual(sim.Snapshot(), control.Snapshot()) {
					t.Fatalf("snapshot diverged at frame %d (vision=%v steering=%v)", frame, vision, steering)
				}
			}
			if len(sim.History()) < 5 || !reflect.DeepEqual(sim.History(), control.History()) {
				t.Fatalf("overlay replay history diverged (vision=%v steering=%v)", vision, steering)
			}
		}
	}
}
