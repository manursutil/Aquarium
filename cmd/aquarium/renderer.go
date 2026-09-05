package main

import (
	"aquarium/internal/simulation"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawFish(f simulation.FishSnapshot) {
	size := f.Genome.Size
	speed := rl.Vector2Length(toRaylibVector(f.Velocity))

	var forward rl.Vector2
	if speed >= 0.001 {
		forward = rl.Vector2Scale(toRaylibVector(f.Velocity), 1/speed)
	} else {
		angle := float64(f.Angle) * math.Pi / 180
		forward = rl.Vector2{
			X: float32(math.Cos(angle)),
			Y: float32(math.Sin(angle)),
		}
	}

	side := rl.Vector2{X: -forward.Y, Y: forward.X}
	point := func(along, across float32) rl.Vector2 {
		return rl.Vector2{
			X: f.Position.X + forward.X*size*along + side.X*size*across,
			Y: f.Position.Y + forward.Y*size*along + side.Y*size*across,
		}
	}

	body := []rl.Vector2{
		toRaylibVector(f.Position),
		point(0.85, 0),
		point(0.45, -0.38),
		point(0, -0.50),
		point(-0.48, -0.34),
		point(-0.58, 0),
		point(-0.48, 0.34),
		point(0, 0.50),
		point(0.45, 0.38),
		point(0.85, 0),
	}
	bodyOutline := body[1:]

	darkColor := rl.ColorBrightness(toRaylibColor(f.Color), -0.25)
	rl.DrawTriangle(
		point(-0.88, -0.48),
		point(-0.48, 0),
		point(-0.88, 0.48),
		darkColor,
	)
	rl.DrawTriangleFan(body, toRaylibColor(f.Color))

	outlineColor := rl.ColorBrightness(toRaylibColor(f.Color), -0.4)
	for i := 0; i < len(bodyOutline)-1; i++ {
		rl.DrawLineEx(bodyOutline[i], bodyOutline[i+1], 1, outlineColor)
	}

	eyePosition := point(0.48, -0.14)
	eyeRadius := max(float32(1.5), size*0.10)
	rl.DrawCircleV(eyePosition, eyeRadius, rl.RayWhite)
	rl.DrawCircleV(eyePosition, eyeRadius*0.45, rl.Black)
}

func toRaylibVector(v simulation.Vector2) rl.Vector2 {
	return rl.Vector2{X: v.X, Y: v.Y}
}

func toRaylibColor(c simulation.Color) rl.Color {
	return rl.NewColor(c.R, c.G, c.B, c.A)
}

func drawAquarium(s simulation.Snapshot) {
	for _, f := range s.Fish {
		drawFish(f)
	}

	for _, f := range s.Food {
		rl.DrawRectangle(int32(f.X), int32(f.Y), 10, 10, rl.Green)
	}

	drawHUD(s)

	mouse := rl.GetMousePosition()
	if i := hoveredFishIndex(s.Fish, mouse); i >= 0 {
		drawFishTooltip(s.Fish[i], mouse, s.MaxHealth)
	}
}
