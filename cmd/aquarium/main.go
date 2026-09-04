package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowWidth      = 800
	WindowHeight     = 600
	WindowTitle      = "Aquarium"
	InitialFishCount = 50
	TargetFPS        = 60
)

func main() {
	aquarium := initAquarium(InitialFishCount)

	rl.InitWindow(WindowWidth, WindowHeight, WindowTitle)
	defer rl.CloseWindow()

	rl.SetTargetFPS(TargetFPS)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)

		dt := rl.GetFrameTime()
		aquarium.update(dt)
		aquarium.draw()

		rl.EndDrawing()
	}
}
