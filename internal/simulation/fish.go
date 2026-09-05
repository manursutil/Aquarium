package simulation

import (
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
)

const (
	NoiseAlpha              = 2.0
	NoiseBeta               = 2.0
	NoiseOctaves            = 3
	NoiseSeed               = 100
	NoiseSampleRate float64 = 0.5
	WanderTurnRate  float32 = 2.0
	MaxHunger       float32 = 1

	FoodSteeringRate       float32 = 10
	AttractionSteeringRate float32 = 15
	repulsionSteeringRate  float32 = 15
	MaxSteeringStrength    float32 = 1
)

var noise = perlin.NewPerlin(NoiseAlpha, NoiseBeta, NoiseOctaves, NoiseSeed)

type Fish struct {
	ID       FishID
	Position Vector2
	Velocity Vector2
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

func newFish(rng *rand.Rand, config Config, genome Genome) Fish {
	return Fish{
		Position: Vector2{
			X: rng.Float32() * config.WorldWidth,
			Y: rng.Float32() * config.WorldHeight,
		},
		Velocity: Vector2{
			X: rng.Float32()*InitialVelocityRange - InitialVelocityOffset,
			Y: rng.Float32()*InitialVelocityRange - InitialVelocityOffset,
		},
		Angle:     rng.Float32() * FullCircleDegrees,
		Genome:    genome,
		NoiseX:    rng.Float64() * InitialNoiseRange,
		Hunger:    0,
		Health:    config.MaxHealth,
		Alive:     true,
		FoodEaten: 0,
		FishEaten: 0,
		Age:       0,
	}
}

func (f *Fish) wrapEdges(config Config) {
	if f.Position.X > config.WorldWidth {
		f.Position.X = 0
	}

	if f.Position.X < 0 {
		f.Position.X = config.WorldWidth
	}

	if f.Position.Y > config.WorldHeight {
		f.Position.Y = 0
	}

	if f.Position.Y < 0 {
		f.Position.Y = config.WorldHeight
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

	f.Velocity = clampMagnitude(f.Velocity, 0, f.Genome.MaxSpeed)
	f.Position = add(f.Position, scale(f.Velocity, dt))
}

func (f *Fish) handleHunger(dt float32, config Config) {
	f.Hunger += hungerRate(*f, config) * dt
	if f.Hunger > MaxHunger {
		f.Hunger = MaxHunger
	}

	if f.Hunger >= MaxHunger {
		f.Health -= config.StarvationDamage * dt
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

func (f *Fish) attractionToSimilarFish(others []Fish, dt float32, config Config) {
	closestIndex := -1
	closestDistance := f.Genome.Vision

	for i := range others {
		if colorDistance(f.Genome.Color, others[i].Genome.Color) < PredationColorDistance {
			dist := distance(f.Position, others[i].Position)
			if dist < closestDistance {
				closestDistance = dist
				closestIndex = i
			}
		}
	}

	if closestIndex == -1 || closestDistance == 0 {
		return
	}

	dir := normalize(subtract(others[closestIndex].Position, f.Position))
	desiredVel := scale(dir, f.Genome.MaxSpeed)

	steeringStrength := f.Genome.AttractionToOthers * AttractionSteeringRate * dt
	if steeringStrength > MaxSteeringStrength {
		steeringStrength = MaxSteeringStrength
	}

	steering := subtract(desiredVel, f.Velocity)
	f.applySteering(add(f.Velocity, scale(steering, steeringStrength)), dt, config)
}

func (f *Fish) steerTowardFood(foods []Food, dt float32, config Config) {
	closestIndex := -1
	closestDistance := f.Genome.Vision

	for i := range foods {
		dist := distance(f.Position, foods[i].Position)
		if dist < closestDistance {
			closestDistance = dist
			closestIndex = i
		}
	}

	if closestIndex == -1 || closestDistance == 0 {
		return
	}

	dir := subtract(foods[closestIndex].Position, f.Position)
	dir = normalize(dir)

	desiredVel := scale(dir, f.Genome.MaxSpeed)

	steeringStrength := f.Genome.AttractionToFood * f.Hunger * FoodSteeringRate * dt
	if steeringStrength > MaxSteeringStrength {
		steeringStrength = MaxSteeringStrength
	}

	steering := subtract(desiredVel, f.Velocity)
	f.applySteering(add(f.Velocity, scale(steering, steeringStrength)), dt, config)
}

func (f *Fish) steerAwayFromPredators(others []Fish, dt float32, config Config) {
	closestIndex := -1
	closestDistance := f.Genome.Vision

	for i := range others {
		if colorDistance(f.Genome.Color, others[i].Genome.Color) > PredationColorDistance {
			dist := distance(f.Position, others[i].Position)
			if dist < closestDistance {
				closestIndex = i
				closestDistance = dist
			}
		}
	}

	if closestIndex == -1 || closestDistance == 0 {
		return
	}

	dir := normalize(subtract(others[closestIndex].Position, f.Position))
	desiredVel := scale(dir, f.Genome.MaxSpeed)

	steeringStrength := f.Genome.FearOfPredators * repulsionSteeringRate * dt
	if steeringStrength > MaxSteeringStrength {
		steeringStrength = MaxSteeringStrength
	}

	steering := subtract(desiredVel, f.Velocity)
	f.applySteering(add(f.Velocity, scale(steering, steeringStrength)), dt, config)
}

func (f *Fish) attackPrey(prey *Fish, config Config) {
	if !prey.Alive {
		return
	}

	wasAlive := prey.Alive
	prey.Health = max(0, prey.Health-CollisionDamage)
	prey.checkDeath()

	f.Health += PredationHealthGain
	if f.Health > config.MaxHealth {
		f.Health = config.MaxHealth
	}

	if wasAlive && !prey.Alive {
		f.FishEaten++
	}
}

func (f *Fish) update(dt float32, config Config) {
	f.Age += dt
	f.move(dt)
	f.wrapEdges(config)
	f.handleHunger(dt, config)
	f.checkDeath()
}

// applySteering bounds each steering response after behavioral weighting.
// A zero BaseAcceleration preserves the unbounded control behavior.
func (f *Fish) applySteering(desired Vector2, dt float32, config Config) {
	if dt <= 0 {
		return
	}
	change := subtract(desired, f.Velocity)
	if config.BaseAcceleration > 0 {
		acceleration := config.BaseAcceleration * (MinFishSize / max(f.Genome.Size, MinFishSize))
		change = clampMagnitude(change, 0, acceleration*dt)
	}
	f.Velocity = clampMagnitude(add(f.Velocity, change), 0, f.Genome.MaxSpeed)
}

func canEat(predator, prey Fish, requiredRatio float32) bool {
	return predator.Genome.Size > prey.Genome.Size &&
		predator.Genome.Size >= prey.Genome.Size*requiredRatio
}
