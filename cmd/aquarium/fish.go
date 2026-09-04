package main

import (
	"math"

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
	rl.DrawCircle(
		int32(f.Position.X),
		int32(f.Position.Y),
		f.Genome.Size,
		f.Genome.Color,
	)
}
