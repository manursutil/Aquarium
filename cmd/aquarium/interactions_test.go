package main

import (
	"aquarium/internal/simulation"
	"reflect"
	"testing"
)

func testApp(t *testing.T) *AppState {
	t.Helper()
	c := simulation.DefaultConfig()
	c.PopulationSize = 8
	c.FoodCount = 5
	c.GenerationDuration = 0.2
	app := &AppState{config: c, seed: 42, view: ViewState{SpeedMultiplier: 1}}
	app.restart(42)
	return app
}

func TestFixedStepBudgetsAndCap(t *testing.T) {
	for _, speed := range []int{1, 5, 20} {
		total := 0
		var remainder float64
		for range 60 {
			n, r := simulationSteps(remainder, fixedStep, speed)
			total += n
			remainder = r
		}
		if total != 60*speed {
			t.Fatalf("%dx: got %d steps", speed, total)
		}
	}
	if n, r := simulationSteps(0, 0, 20); n != 0 || r != 0 {
		t.Fatal(n, r)
	}
	if n, r := simulationSteps(0, 100, 20); n != maxStepsPerFrame || r != 0 {
		t.Fatal(n, r)
	}
	app := testApp(t)
	before := app.sim.Snapshot()
	app.view.Paused = true
	app.accumulator = float64(fixedStep) / 2
	app.advance(10)
	if !reflect.DeepEqual(before, app.sim.Snapshot()) || app.accumulator != float64(fixedStep)/2 {
		t.Fatal("pause advanced time")
	}
}

func TestRestartAndControlledReplay(t *testing.T) {
	app := testApp(t)
	initial := app.sim.Snapshot()
	for range 120 {
		app.advance(fixedStep)
	}
	history := app.sim.History()
	app.view.HasSelectedFish = true
	app.view.LastHistoryCount = 99
	app.view.SummaryVisibleFor = 2
	app.accumulator = 0.01
	app.restart(42)
	if !reflect.DeepEqual(initial, app.sim.Snapshot()) || app.view.HasSelectedFish || app.view.LastHistoryCount != 0 || app.view.SummaryVisibleFor != 0 || app.accumulator != 0 {
		t.Fatal("restart did not reset run")
	}
	for range 120 {
		app.advance(fixedStep)
	}
	if !reflect.DeepEqual(history, app.sim.History()) {
		t.Fatal("restart replay diverged")
	}
	control := testApp(t)
	for range 120 {
		control.sim.Step(fixedStep)
	}
	if !reflect.DeepEqual(control.sim.History(), history) {
		t.Fatal("controls changed history")
	}
	app.restart(43)
	if app.seed != 43 || reflect.DeepEqual(initial, app.sim.Snapshot()) || len(app.sim.Snapshot().Fish) != app.config.PopulationSize {
		t.Fatal("new seed/config")
	}
}

func TestSelectionCompactionAndRollover(t *testing.T) {
	view := ViewState{HasSelectedFish: true, SelectedFishID: 7}
	fish := []simulation.FishSnapshot{{ID: 1}, {ID: 7}}
	view.clearMissingSelection(fish[1:])
	if !view.HasSelectedFish {
		t.Fatal("compaction lost selection")
	}
	view.clearMissingSelection([]simulation.FishSnapshot{{ID: 8}})
	if view.HasSelectedFish {
		t.Fatal("missing selection persisted")
	}
	app := testApp(t)
	app.view.HasSelectedFish = true
	app.view.SelectedFishID = app.sim.Snapshot().Fish[0].ID
	for range 30 {
		app.advance(fixedStep)
	}
	app.view.clearMissingSelection(app.sim.Snapshot().Fish)
	if app.view.HasSelectedFish {
		t.Fatal("rollover retained selection")
	}
}

func TestSummaryTransitions(t *testing.T) {
	view := ViewState{}
	view.updateSummary(0, 0.1)
	if view.SummaryVisibleFor != 0 {
		t.Fatal("empty history activated summary")
	}
	view.updateSummary(1, 0.1)
	if view.SummaryVisibleFor != 2 {
		t.Fatal("first record")
	}
	view.updateSummary(1, 0.5)
	if view.SummaryVisibleFor != 1.5 {
		t.Fatal("same record reactivated")
	}
	view.updateSummary(4, 0.1)
	if view.LastHistoryCount != 4 || view.SummaryVisibleFor != 2 {
		t.Fatal("multiple records")
	}
	view.updateSummary(4, 3)
	if view.SummaryVisibleFor != 0 {
		t.Fatal("summary did not expire")
	}
}
