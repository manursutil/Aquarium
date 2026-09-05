package simulation

import (
	"fmt"
	"math"
)

// Config contains the tunable world and evolution parameters.
type Config struct {
	WorldWidth, WorldHeight                                      float32
	PopulationSize, FoodCount                                    int
	GenerationDuration                                           float32
	MaxHealth, FoodHealing, StarvationDamage, PredationSizeRatio float32
	MutationRate, MutationSize                                   float32
	EliteCount                                                   int

	SpeedEnergyWeight  float32
	SizeEnergyWeight   float32
	VisionEnergyWeight float32
	BaseAcceleration   float32
}

func DefaultConfig() Config {
	return Config{
		WorldWidth:         800,
		WorldHeight:        600,
		PopulationSize:     50,
		FoodCount:          35,
		GenerationDuration: 30,
		MaxHealth:          10,
		FoodHealing:        3,
		StarvationDamage:   1,
		PredationSizeRatio: 1,
		MutationRate:       0.15,
		MutationSize:       0.10,
		EliteCount:         1,
		SpeedEnergyWeight:  0.20,
		SizeEnergyWeight:   0.10,
		VisionEnergyWeight: 0.10,
		BaseAcceleration:   30,
	}
}

func (c Config) Validate() error {
	for name, value := range map[string]float32{
		"world width":          c.WorldWidth,
		"world height":         c.WorldHeight,
		"generation duration":  c.GenerationDuration,
		"max health":           c.MaxHealth,
		"predation size ratio": c.PredationSizeRatio,
	} {
		if value <= 0 || math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return fmt.Errorf("%s must be finite and positive", name)
		}
	}

	for name, value := range map[string]float32{
		"food healing":      c.FoodHealing,
		"starvation damage": c.StarvationDamage,
		"mutation rate":     c.MutationRate,
		"mutation size":     c.MutationSize,
	} {
		if value < 0 || math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return fmt.Errorf("%s must be finite and nonnegative", name)
		}
	}

	for name, value := range map[string]float32{
		"speed energy weight":  c.SpeedEnergyWeight,
		"size energy weight":   c.SizeEnergyWeight,
		"vision energy weight": c.VisionEnergyWeight,
		"base acceleration":    c.BaseAcceleration,
	} {
		if value < 0 || math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return fmt.Errorf("%s must be finite and nonnegative", name)
		}
	}

	if c.PopulationSize <= 0 {
		return fmt.Errorf("population size must be positive")
	}

	if c.FoodCount < 0 {
		return fmt.Errorf("food count must be nonnegative")
	}

	if c.EliteCount < 0 || c.EliteCount > c.PopulationSize {
		return fmt.Errorf("elite count is outside population bounds")
	}

	if c.MutationRate > 1 || c.MutationSize > 1 {
		return fmt.Errorf("mutation rate and size must not exceed one")
	}

	if c.PredationSizeRatio < 1 {
		return fmt.Errorf("predation size ratio must be at least one")
	}

	return nil
}
