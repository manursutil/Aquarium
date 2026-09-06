package simulation

import (
	"math"
	"math/rand"
	"slices"
)

const (
	CollisionDamage        float32 = 4
	InitialVelocityRange   float32 = 200
	InitialVelocityOffset  float32 = 100
	FullCircleDegrees      float32 = 360
	InitialNoiseRange      float64 = 1000
	FoodHungerReduction    float32 = 0.3
	PredationColorDistance float32 = 100
	PredationHealthGain    float32 = 1

	FoodSize           = 10
	FoodRadius float32 = FoodSize / 2
)

type Simulation struct {
	config  Config
	lineage map[FishID]LineageRecord
	rng     *rand.Rand
	nextID  FishID
	fish    []Fish
	food    []Food

	generation           int
	generationElapsed    float32
	generationCandidates []Candidate
	history              []GenerationStats
}

func New(config Config, seed int64) (*Simulation, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	a := &Simulation{config: config, rng: rand.New(rand.NewSource(seed)), generation: 1, lineage: make(map[FishID]LineageRecord)}

	a.fish = make([]Fish, config.PopulationSize)
	for i := range a.fish {
		a.fish[i] = a.spawnOffspring(Offspring{Genome: initRandomGenome(a.rng)}, 1)
	}

	a.food = make([]Food, config.FoodCount)
	for i := range a.food {
		a.food[i] = Food{HealAmount: config.FoodHealing}
		a.food[i].reposition(a.rng, config)
	}

	return a, nil
}

func (a *Simulation) spawnOffspring(child Offspring, bornIn int) Fish {
	genome := child.Genome
	f := newFish(a.rng, a.config, genome)

	a.nextID++
	f.ID = a.nextID
	f.ParentA, f.ParentB, f.BornIn, f.Elite = child.ParentA, child.ParentB, bornIn, child.Elite
	if a.lineage == nil {
		a.lineage = make(map[FishID]LineageRecord)
	}
	a.lineage[f.ID] = recordFromFish(f)

	return f
}

func (a *Simulation) removeDead() {
	alive := a.fish[:0]

	for _, f := range a.fish {
		if f.Alive {
			alive = append(alive, f)
			continue
		}

		a.generationCandidates = append(a.generationCandidates, candidateFromFish(f))
		a.finalizeLineage(f, "death")
	}

	a.fish = alive
}

func (a *Simulation) Step(dt float32) {
	for i := range len(a.fish) {
		a.fish[i].steering = SteeringSnapshot{}
		before := a.fish[i].Velocity

		a.fish[i].attractionToSimilarFish(a.fish, dt, a.config)
		a.fish[i].steerTowardFood(a.food, dt, a.config)
		a.fish[i].steerAwayFromPredators(a.fish, dt, a.config)

		a.fish[i].steering.Final = subtract(a.fish[i].Velocity, before)

		a.fish[i].update(dt, a.config)
	}
	a.handleFishFoodCollisions()
	a.handleFishFishCollisions()
	a.removeDead()

	a.generationElapsed += dt

	if a.generationElapsed >= a.config.GenerationDuration {
		a.generationElapsed -= a.config.GenerationDuration
		a.advanceGeneration()
	} else if len(a.fish) == 0 {
		a.generationElapsed = 0
		a.advanceGeneration()
	}
}

func (a *Simulation) handleFishFoodCollisions() {
	for i := range a.fish {
		f := &a.fish[i]
		if !f.Alive {
			continue
		}

		for j := range slices.Backward(a.food) {
			food := &a.food[j]
			dist := distance(f.Position, food.Position)
			if dist < f.Genome.Size+FoodRadius {
				f.Health += food.HealAmount
				if f.Health > a.config.MaxHealth {
					f.Health = a.config.MaxHealth
				}

				f.Hunger -= FoodHungerReduction
				if f.Hunger < 0 {
					f.Hunger = 0
				}
				f.FoodEaten++
				food.reposition(a.rng, a.config)
			}
		}
	}
}

func colorDistance(c1 Color, c2 Color) float32 {
	dr := float32(c1.R) - float32(c2.R)
	dg := float32(c1.G) - float32(c2.G)
	db := float32(c1.B) - float32(c2.B)

	return float32(math.Sqrt(float64(dr*dr + dg*dg + db*db)))
}

func (a *Simulation) handleFishFishCollisions() {
	for i := range a.fish {
		if !a.fish[i].Alive {
			continue
		}

		for j := i + 1; j < len(a.fish); j++ {
			if !a.fish[j].Alive {
				continue
			}

			f1, f2 := &a.fish[i], &a.fish[j]
			delta := subtract(f2.Position, f1.Position)
			dist := length(delta)
			minDist := f1.Genome.Size + f2.Genome.Size
			if dist < minDist && dist > 0 {

				// Physical separation (avoid overlapping)
				overlap := minDist - dist
				push := scale(normalize(delta), overlap/2)
				f1.Position = subtract(f1.Position, push)
				f2.Position = add(f2.Position, push)

				// Eating mechanic
				colorDist := colorDistance(f1.Genome.Color, f2.Genome.Color)

				if colorDist > PredationColorDistance {
					if canEat(*f1, *f2, a.config.PredationSizeRatio) {
						f1.attackPrey(f2, a.config)
					} else if canEat(*f2, *f1, a.config.PredationSizeRatio) {
						f2.attackPrey(f1, a.config)
						if !f1.Alive {
							break
						}
					}
				}
			}
		}
	}
}

func (a *Simulation) advanceGeneration() {
	candidates := append([]Candidate{}, a.generationCandidates...)
	candidates = append(candidates, makeCandidates(a.fish)...)

	for _, f := range a.fish {
		a.finalizeLineage(f, "rollover")
	}

	stats := summarizeGeneration(a.generation, len(a.fish), candidates)
	a.history = append(a.history, stats)

	nextGenomes := evolve(a.rng, candidates, a.config)

	nextFish := make([]Fish, len(nextGenomes))
	for i, genome := range nextGenomes {
		nextFish[i] = a.spawnOffspring(genome, a.generation+1)
	}

	a.fish = nextFish
	a.generationCandidates = nil
	a.generation++
	a.evictLineage()
}

func (a *Simulation) bestCurrentFitness() float32 {
	best := float32(0)

	for _, candidate := range a.generationCandidates {
		if candidate.Fitness > best {
			best = candidate.Fitness
		}
	}

	for _, fish := range a.fish {
		score := fitness(fish)
		if score > best {
			best = score
		}
	}

	return best
}

func (a *Simulation) History() []GenerationStats {
	return append([]GenerationStats(nil), a.history...)
}
