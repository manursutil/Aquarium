package main

import (
	"log"
	"math/rand/v2"

	"aquarium/internal/simulation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowTitle              = "Aquarium"
	TargetFPS                = 60
	fixedStep        float32 = 1.0 / 60.0
	maxStepsPerFrame         = 40
)

type AppState struct {
	config      simulation.Config
	seed        int64
	sim         *simulation.Simulation
	view        ViewState
	accumulator float64
}

func main() {
	app := AppState{config: simulation.DefaultConfig(), seed: 42, view: ViewState{SpeedMultiplier: 1}}
	app.restart(app.seed)
	config := app.config

	rl.InitWindow(int32(config.WorldWidth), int32(config.WorldHeight), WindowTitle)
	defer rl.CloseWindow()
	rl.SetExitKey(rl.KeyNull)

	rl.SetTargetFPS(TargetFPS)

	for !rl.WindowShouldClose() {
		updateViewState(&app.view)
		if rl.IsKeyPressed(rl.KeyR) {
			app.restart(app.seed)
		}
		if rl.IsKeyPressed(rl.KeyN) {
			seed := rand.Int64()
			for seed == app.seed {
				seed = rand.Int64()
			}
			app.restart(seed)
		}
		frameTime := rl.GetFrameTime()
		app.advance(frameTime)
		snapshot := app.sim.Snapshot()
		history := app.sim.History()
		pinFish(&app.view, snapshot.Fish, rl.GetMousePosition())
		app.view.clearMissingSelection(snapshot.Fish)
		app.view.updateSummary(len(history), frameTime)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)

		drawAquarium(snapshot, history, app.view, app.seed)

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

func (app *AppState) restart(seed int64) {
	app.sim = NewSimulation(app.config, seed)
	app.seed = seed
	app.accumulator = 0
	app.view.HasSelectedFish = false
	app.view.SelectedFishID = 0
	app.view.LastHistoryCount = 0
	app.view.SummaryVisibleFor = 0
}

// Use float64 for the clock so 20 steps do not round down to 19.
func simulationSteps(accumulator float64, frameTime float32, speed int) (int, float64) {
	budget := accumulator + float64(frameTime)*float64(speed)
	steps := int(budget / float64(fixedStep))
	if steps >= maxStepsPerFrame {
		return maxStepsPerFrame, 0
	}
	return steps, budget - float64(steps)*float64(fixedStep)
}

func (app *AppState) advance(frameTime float32) {
	if app.view.Paused {
		return
	}
	steps, remainder := simulationSteps(app.accumulator, frameTime, app.view.SpeedMultiplier)
	app.accumulator = remainder
	for range steps {
		app.sim.Step(fixedStep)
	}
}
