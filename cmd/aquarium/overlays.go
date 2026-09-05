package main

import (
	"math"

	"aquarium/internal/simulation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type overlayVector struct {
	Vector    simulation.Vector2
	Color     rl.Color
	Thickness float32
}

type fishOverlay struct {
	Fish       simulation.FishSnapshot
	ShowVision bool
	Vectors    []overlayVector
}

func buildFishOverlay(fish []simulation.FishSnapshot, view ViewState) fishOverlay {
	if !view.HasSelectedFish || (!view.ShowVision && !view.ShowSteering) {
		return fishOverlay{}
	}
	f, found := fishByID(fish, view.SelectedFishID)
	if !found {
		return fishOverlay{}
	}
	overlay := fishOverlay{Fish: f, ShowVision: view.ShowVision}
	if view.ShowSteering {
		// Draw the heavier final response underneath the component lines.
		overlay.Vectors = []overlayVector{
			{displayVector(f.Steering.Final, 80), rl.White, 4},
			{displayVector(f.Steering.Food, 60), rl.Lime, 2},
			{displayVector(f.Steering.Social, 60), rl.Gold, 2},
			{displayVector(f.Steering.Predator, 60), rl.Magenta, 2},
		}
	}
	return overlay
}

func displayVector(v simulation.Vector2, pixels float32) simulation.Vector2 {
	length := math.Hypot(float64(v.X), float64(v.Y))
	if length == 0 {
		return simulation.Vector2{}
	}
	return simulation.Vector2{
		X: float32(float64(v.X) / length * float64(pixels)),
		Y: float32(float64(v.Y) / length * float64(pixels)),
	}
}
