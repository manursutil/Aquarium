package main

import (
	"aquarium/internal/simulation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ViewState struct {
	Paused          bool
	SpeedMultiplier int
	ShowVision      bool
	ShowSteering    bool
	SelectedFishID  simulation.FishID
	HasSelectedFish bool

	SummaryVisibleFor float32
	LastHistoryCount  int
}

func updateViewState(view *ViewState) {
	if rl.IsKeyPressed(rl.KeySpace) {
		view.Paused = !view.Paused
	}
	if rl.IsKeyPressed(rl.KeyOne) {
		view.SpeedMultiplier = 1
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		view.SpeedMultiplier = 5
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		view.SpeedMultiplier = 20
	}
	if rl.IsKeyPressed(rl.KeyV) {
		view.ShowVision = !view.ShowVision
	}
	if rl.IsKeyPressed(rl.KeyF) {
		view.ShowSteering = !view.ShowSteering
	}
	if rl.IsKeyPressed(rl.KeyEscape) {
		view.HasSelectedFish = false
	}
}
