package main

import (
	"aquarium/internal/simulation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ViewState struct {
	AncestryID        simulation.FishID
	HideFastSummaries bool
	Paused            bool
	SpeedMultiplier   int
	ShowVision        bool
	ShowSteering      bool
	SelectedFishID    simulation.FishID
	HasSelectedFish   bool

	SummaryVisibleFor float32
	LastHistoryCount  int
}

func updateViewState(view *ViewState) {
	if rl.IsKeyPressed(rl.KeyS) {
		view.HideFastSummaries = !view.HideFastSummaries
	}
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
		view.AncestryID = 0
	}
}

func (view *ViewState) clearMissingSelection(fish []simulation.FishSnapshot) {
	if _, found := fishByID(fish, view.SelectedFishID); view.HasSelectedFish && !found {
		view.HasSelectedFish = false
	}
}

func (view *ViewState) updateSummary(count int, frameTime float32) {
	view.SummaryVisibleFor = max(0, view.SummaryVisibleFor-frameTime)
	if count > view.LastHistoryCount {
		view.SummaryVisibleFor = 2
	}
	view.LastHistoryCount = count
}
