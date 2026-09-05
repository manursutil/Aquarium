package main

import (
	"log"

	"aquarium/internal/simulation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowTitle              = "Aquarium"
	TargetFPS                = 60
	fixedStep        float32 = 1.0 / 60.0
	maxStepsPerFrame         = 10
)

type AppState struct {
	config simulation.Config
	seed   int64
	sim    *simulation.Simulation
	view   ViewState
}

func main() {
	config := simulation.DefaultConfig()
	aquarium, err := simulation.New(config, 42)
	if err != nil {
		log.Fatal(err)
	}

	view := ViewState{SpeedMultiplier: 1}

	rl.InitWindow(int32(config.WorldWidth), int32(config.WorldHeight), WindowTitle)
	defer rl.CloseWindow()
	rl.SetExitKey(rl.KeyNull)

	rl.SetTargetFPS(TargetFPS)

	var accumulator float32

	for !rl.WindowShouldClose() {
		updateViewState(&view)

		frameTime := rl.GetFrameTime()
		if !view.Paused {
			accumulator += frameTime * float32(view.SpeedMultiplier)
		}

		steps := 0
		for accumulator >= fixedStep && steps < maxStepsPerFrame {
			aquarium.Step(fixedStep)
			accumulator -= fixedStep
			steps++
		}

		if steps == maxStepsPerFrame && accumulator >= fixedStep {
			accumulator = 0
		}

		snapshot := aquarium.Snapshot()
		history := aquarium.History()

		pinFish(&view, snapshot.Fish, rl.GetMousePosition())
		if view.HasSelectedFish {
			if _, found := fishByID(snapshot.Fish, view.SelectedFishID); !found {
				view.HasSelectedFish = false
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)

		drawAquarium(snapshot, history, view)

		rl.EndDrawing()
	}
}

func NewSimulation(config simulation.Config, seed int64) *simulation.Simulation {
	sim, err := simulation.New(config, seed)
	if err != nil {
		log.Fatal(err)
	}

	return sim
}

func pinFish(view *ViewState, fish []simulation.FishSnapshot, mouse rl.Vector2) {
	if !rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		return
	}

	if index := hoveredFishIndex(fish, mouse); index >= 0 {
		view.SelectedFishID = fish[index].ID
		view.HasSelectedFish = true
	}
}
