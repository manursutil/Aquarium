package main

import (
	"aquarium/internal/simulation"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowTitle = "Aquarium"
	TargetFPS   = 60
)

func main() {
	config := simulation.DefaultConfig()
	aquarium, err := simulation.New(config, 42)
	if err != nil {
		log.Fatal(err)
	}

	rl.InitWindow(int32(config.WorldWidth), int32(config.WorldHeight), WindowTitle)
	defer rl.CloseWindow()

	rl.SetTargetFPS(TargetFPS)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		aquarium.Step(dt)
		snapshot := aquarium.Snapshot()
		history := aquarium.History()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)
		drawAquarium(snapshot, history)

		rl.EndDrawing()
	}
}
