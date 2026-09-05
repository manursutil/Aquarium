package simulation

type FishID uint64

// FishSnapshot is a value-only display record, independent of live world state.
type FishSnapshot struct {
	ID                           FishID
	Position, Velocity           Vector2
	Angle, Radius                float32
	Color                        Color
	Health, Hunger, Age, Fitness float32
	FoodEaten, FishEaten         int
	Genome                       Genome
}

// Snapshot owns its slices; callers may retain or modify them safely.
type Snapshot struct {
	Generation                            int
	GenerationElapsed, GenerationDuration float32
	MaxHealth                             float32
	Evaluated                             int
	BestFitness                           float32
	Fish                                  []FishSnapshot
	Food                                  []Vector2
}

func (s *Simulation) Snapshot() Snapshot {
	fish := make([]FishSnapshot, len(s.fish))
	for i, f := range s.fish {
		fish[i] = FishSnapshot{ID: f.ID, Position: f.Position, Velocity: f.Velocity, Angle: f.Angle,
			Radius: f.Genome.Size, Color: f.Genome.Color, Health: f.Health, Hunger: f.Hunger, Age: f.Age,
			Fitness: fitness(f), FoodEaten: f.FoodEaten, FishEaten: f.FishEaten, Genome: f.Genome}
	}
	food := make([]Vector2, len(s.food))
	for i, f := range s.food {
		food[i] = f.Position
	}
	return Snapshot{Generation: s.generation, GenerationElapsed: s.generationElapsed,
		GenerationDuration: s.config.GenerationDuration, MaxHealth: s.config.MaxHealth,
		Evaluated: len(s.generationCandidates), BestFitness: s.bestCurrentFitness(), Fish: fish, Food: food}
}
