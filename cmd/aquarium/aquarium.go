package main

import (
	"math"
	"math/rand"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MaxHealth              float32 = 10
	FoodHealing            float32 = 3
	CollisionDamage        float32 = 4
	InitialVelocityRange   float32 = 200
	InitialVelocityOffset  float32 = 100
	FullCircleDegrees      float32 = 360
	InitialNoiseRange      float64 = 1000
	FoodHungerReduction    float32 = 0.3
	PredationColorDistance float32 = 100
	PredationHealthGain    float32 = 1
	FishPerFood                    = 2
	BaseFoodCount                  = 10
	FoodSize                       = 10
	FoodRadius             float32 = FoodSize / 2
	GenerationDuration     float32 = 30
)

type Aquarium struct {
	Fish []Fish
	Food []Food

	generationElapsed    float32
	generationCandidates []Candidate
}

func initAquarium(n int) Aquarium {
	fish := make([]Fish, n)

	for i := range n {
		fish[i] = newFish(initRandomGenome())
	}

	foodNumber := n/FishPerFood + BaseFoodCount
	food := make([]Food, foodNumber)

	for i := range foodNumber {
		food[i] = Food{
			Position: rl.Vector2{
				X: rand.Float32() * WindowWidth,
				Y: rand.Float32() * WindowHeight,
			},
			HealAmount: FoodHealing,
		}
	}

	return Aquarium{Fish: fish, Food: food}
}

func (a *Aquarium) removeDead() {
	alive := a.Fish[:0]

	for _, f := range a.Fish {
		if f.Alive {
			alive = append(alive, f)
			continue
		}

		a.generationCandidates = append(a.generationCandidates, Candidate{
			Genome:  f.Genome,
			Fitness: fitness(f),
		})
	}

	a.Fish = alive
}

func (a *Aquarium) update(dt float32) {
	for i := range len(a.Fish) {
		a.Fish[i].attractionToSimilarFish(a.Fish, dt)
		a.Fish[i].steerTowardFood(a.Food, dt)
		a.Fish[i].steerAwayFromPredators(a.Fish, dt)
		a.Fish[i].update(dt)
	}
	a.handleFishFoodCollisions()
	a.handleFishFishCollisions()
	a.removeDead()

	a.generationElapsed += dt

	if a.generationElapsed >= GenerationDuration || len(a.Fish) == 0 {
		a.advanceGeneration()
	}
}

func (a *Aquarium) handleFishFoodCollisions() {
	for i := range a.Fish {
		f := &a.Fish[i]
		if !f.Alive {
			continue
		}

		for j := range slices.Backward(a.Food) {
			food := &a.Food[j]
			dist := rl.Vector2Distance(f.Position, food.Position)
			if dist < f.Genome.Size+FoodRadius {
				f.Health += food.HealAmount
				if f.Health > MaxHealth {
					f.Health = MaxHealth
				}

				f.Hunger -= FoodHungerReduction
				if f.Hunger < 0 {
					f.Hunger = 0
				}
				f.FoodEaten++
				food.reposition()
			}
		}
	}
}

func colorDistance(c1 rl.Color, c2 rl.Color) float32 {
	dr := float32(c1.R) - float32(c2.R)
	dg := float32(c1.G) - float32(c2.G)
	db := float32(c1.B) - float32(c2.B)

	return float32(math.Sqrt(float64(dr*dr + dg*dg + db*db)))
}

func (a *Aquarium) handleFishFishCollisions() {
	for i := range a.Fish {
		if !a.Fish[i].Alive {
			continue
		}

		for j := i + 1; j < len(a.Fish); j++ {
			if !a.Fish[j].Alive {
				continue
			}

			f1, f2 := &a.Fish[i], &a.Fish[j]
			delta := rl.Vector2Subtract(f2.Position, f1.Position)
			dist := rl.Vector2Length(delta)
			minDist := f1.Genome.Size + f2.Genome.Size
			if dist < minDist && dist > 0 {

				// Physical separation (avoid overlapping)
				overlap := minDist - dist
				push := rl.Vector2Scale(rl.Vector2Normalize(delta), overlap/2)
				f1.Position = rl.Vector2Subtract(f1.Position, push)
				f2.Position = rl.Vector2Add(f2.Position, push)

				// Eating mechanic
				colorDist := colorDistance(f1.Genome.Color, f2.Genome.Color)

				if colorDist > PredationColorDistance {
					if f1.Genome.Size > f2.Genome.Size {
						f1.attackPrey(f2)
					} else if f1.Genome.Size < f2.Genome.Size {
						f2.attackPrey(f1)
					}
				}
			}
		}
	}
}

func (a *Aquarium) advanceGeneration() {
	candidates := append([]Candidate{}, a.generationCandidates...)
	candidates = append(candidates, makeCandidates(a.Fish)...)

	nextGenomes := evolve(candidates, InitialFishCount)

	nextFish := make([]Fish, len(nextGenomes))
	for i, genome := range nextGenomes {
		nextFish[i] = newFish(genome)
	}

	a.Fish = nextFish
	a.generationCandidates = nil
	a.generationElapsed = 0
}

func (a *Aquarium) draw() {
	for i := range len(a.Fish) {
		if a.Fish[i].Alive {
			a.Fish[i].draw()
		}
	}

	for i := range len(a.Food) {
		f := a.Food[i]
		rl.DrawRectangle(int32(f.Position.X), int32(f.Position.Y), FoodSize, FoodSize, rl.Green)
	}
}
