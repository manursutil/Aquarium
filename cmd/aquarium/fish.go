package main

import (
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	NoiseAlpha                        = 2.0
	NoiseBeta                         = 2.0
	NoiseOctaves                      = 3
	NoiseSeed                         = 100
	NoiseSampleRate           float64 = 0.5
	WanderTurnRate            float32 = 2.0
	MaxHunger                 float32 = 1
	StarvationDamagePerSecond float32 = 1
	FoodSteeringRate          float32 = 10
	AttractionSteeringRate    float32 = 15
	repulsionSteeringRate     float32 = 15
	MaxSteeringStrength       float32 = 1
)

var noise = perlin.NewPerlin(NoiseAlpha, NoiseBeta, NoiseOctaves, NoiseSeed)

type Fish struct {
	Position rl.Vector2
	Velocity rl.Vector2
	Angle    float32
	Genome   Genome

	NoiseX float64

	Hunger float32
	Health float32
	Alive  bool

	// Performace metrics
	FoodEaten int
	FishEaten int
	Age       float32
}

func newFish(genome Genome) Fish {
	return Fish{
		Position: rl.Vector2{
			X: rand.Float32() * WindowWidth,
			Y: rand.Float32() * WindowHeight,
		},
		Velocity: rl.Vector2{
			X: rand.Float32()*InitialVelocityRange - InitialVelocityOffset,
			Y: rand.Float32()*InitialVelocityRange - InitialVelocityOffset,
		},
		Angle:     rand.Float32() * FullCircleDegrees,
		Genome:    genome,
		NoiseX:    rand.Float64() * InitialNoiseRange,
		Hunger:    0,
		Health:    MaxHealth,
		Alive:     true,
		FoodEaten: 0,
		FishEaten: 0,
		Age:       0,
	}
}

func (f *Fish) wrapEdges() {
	if f.Position.X > WindowWidth {
		f.Position.X = 0
	}

	if f.Position.X < 0 {
		f.Position.X = WindowWidth
	}

	if f.Position.Y > WindowHeight {
		f.Position.Y = 0
	}

	if f.Position.Y < 0 {
		f.Position.Y = WindowHeight
	}
}

func (f *Fish) move(dt float32) {
	n := noise.Noise1D(f.NoiseX)
	f.NoiseX += NoiseSampleRate * float64(dt)

	turnAmount := float32(n) * WanderTurnRate * dt
	cos := float32(math.Cos(float64(turnAmount)))
	sin := float32(math.Sin(float64(turnAmount)))

	vx := f.Velocity.X
	vy := f.Velocity.Y

	f.Velocity.X = vx*cos - vy*sin
	f.Velocity.Y = vx*sin + vy*cos

	f.Velocity = rl.Vector2ClampValue(f.Velocity, 0, f.Genome.MaxSpeed)
	f.Position = rl.Vector2Add(f.Position, rl.Vector2Scale(f.Velocity, dt))
}

func (f *Fish) handleHunger(dt float32) {
	f.Hunger += f.Genome.Metabolism * dt
	if f.Hunger > MaxHunger {
		f.Hunger = MaxHunger
	}

	if f.Hunger >= MaxHunger {
		f.Health -= StarvationDamagePerSecond * dt
		if f.Health < 0 {
			f.Health = 0
		}
	}
}

func (f *Fish) checkDeath() {
	if f.Health <= 0 {
		f.Alive = false
	}
}

func (f *Fish) attractionToSimilarFish(others []Fish, dt float32) {
	closestIndex := -1
	closestDistance := f.Genome.Vision

	for i := range others {
		if colorDistance(f.Genome.Color, others[i].Genome.Color) < PredationColorDistance {
			dist := rl.Vector2Distance(f.Position, others[i].Position)
			if dist < closestDistance {
				closestDistance = dist
				closestIndex = i
			}
		}
	}

	if closestIndex == -1 || closestDistance == 0 {
		return
	}

	dir := rl.Vector2Normalize(rl.Vector2Subtract(others[closestIndex].Position, f.Position))
	desiredVel := rl.Vector2Scale(dir, f.Genome.MaxSpeed)

	steeringStrength := f.Genome.AttractionToOthers * AttractionSteeringRate * dt
	if steeringStrength > MaxSteeringStrength {
		steeringStrength = MaxSteeringStrength
	}

	steering := rl.Vector2Subtract(desiredVel, f.Velocity)
	f.Velocity = rl.Vector2Add(f.Velocity, rl.Vector2Scale(steering, steeringStrength))
}

func (f *Fish) steerTowardFood(foods []Food, dt float32) {
	closestIndex := -1
	closestDistance := f.Genome.Vision

	for i := range foods {
		dist := rl.Vector2Distance(f.Position, foods[i].Position)
		if dist < closestDistance {
			closestDistance = dist
			closestIndex = i
		}
	}

	if closestIndex == -1 || closestDistance == 0 {
		return
	}

	dir := rl.Vector2Subtract(foods[closestIndex].Position, f.Position)
	dir = rl.Vector2Normalize(dir)

	desiredVel := rl.Vector2Scale(dir, f.Genome.MaxSpeed)

	steeringStrength := f.Genome.AttractionToFood * f.Hunger * FoodSteeringRate * dt
	if steeringStrength > MaxSteeringStrength {
		steeringStrength = MaxSteeringStrength
	}

	steering := rl.Vector2Subtract(desiredVel, f.Velocity)
	f.Velocity = rl.Vector2Add(f.Velocity, rl.Vector2Scale(steering, steeringStrength))
}

func (f *Fish) steerAwayFromPredators(others []Fish, dt float32) {
	closestIndex := -1
	closestDistance := f.Genome.Vision

	for i := range others {
		if colorDistance(f.Genome.Color, others[i].Genome.Color) > PredationColorDistance {
			dist := rl.Vector2Distance(f.Position, others[i].Position)
			if dist < closestDistance {
				closestIndex = i
				closestDistance = dist
			}
		}
	}

	if closestIndex == -1 || closestDistance == 0 {
		return
	}

	dir := rl.Vector2Normalize(rl.Vector2Subtract(others[closestIndex].Position, f.Position))
	desiredVel := rl.Vector2Scale(dir, f.Genome.MaxSpeed)

	steeringStrength := f.Genome.FearOfPredators * repulsionSteeringRate * dt
	if steeringStrength > MaxSteeringStrength {
		steeringStrength = MaxSteeringStrength
	}

	steering := rl.Vector2Subtract(desiredVel, f.Velocity)
	f.Velocity = rl.Vector2Add(f.Velocity, rl.Vector2Scale(steering, steeringStrength))
}

func (f *Fish) attackPrey(prey *Fish) {
	if !prey.Alive {
		return
	}

	wasAlive := prey.Alive
	prey.Health -= CollisionDamage
	prey.checkDeath()

	f.Health += PredationHealthGain
	if f.Health > MaxHealth {
		f.Health = MaxHealth
	}

	if wasAlive && !prey.Alive {
		f.FishEaten++
	}
}

func (f *Fish) update(dt float32) {
	f.Age += dt
	f.move(dt)
	f.wrapEdges()
	f.handleHunger(dt)
	f.checkDeath()
}

func (f *Fish) draw() {
	size := f.Genome.Size
	speed := rl.Vector2Length(f.Velocity)
	var forward rl.Vector2
	if speed >= 0.001 {
		forward = rl.Vector2Scale(f.Velocity, 1/speed)
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
		f.Position,
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

	darkColor := rl.ColorBrightness(f.Genome.Color, -0.25)
	rl.DrawTriangle(
		point(-0.88, -0.48),
		point(-0.48, 0),
		point(-0.88, 0.48),
		darkColor,
	)
	rl.DrawTriangleFan(body, f.Genome.Color)

	outlineColor := rl.ColorBrightness(f.Genome.Color, -0.4)
	for i := 0; i < len(bodyOutline)-1; i++ {
		rl.DrawLineEx(bodyOutline[i], bodyOutline[i+1], 1, outlineColor)
	}

	eyePosition := point(0.48, -0.14)
	eyeRadius := max(float32(1.5), size*0.10)
	rl.DrawCircleV(eyePosition, eyeRadius, rl.RayWhite)
	rl.DrawCircleV(eyePosition, eyeRadius*0.45, rl.Black)
}
