package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Food struct {
	Position   rl.Vector2
	HealAmount float32
}

func (f *Food) reposition() {
	f.Position = rl.Vector2{
		X: rand.Float32() * WindowWidth,
		Y: rand.Float32() * WindowHeight,
	}
}
