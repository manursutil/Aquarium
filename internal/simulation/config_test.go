package simulation_test

import (
	"aquarium/internal/simulation"
	"math"
	"testing"
)

func TestNewRejectsInvalidConfiguration(t *testing.T) {
	cases := []func(*simulation.Config){
		func(c *simulation.Config) { c.WorldWidth = 0 },
		func(c *simulation.Config) { c.WorldHeight = float32(math.NaN()) },
		func(c *simulation.Config) { c.GenerationDuration = float32(math.Inf(1)) },
		func(c *simulation.Config) { c.PopulationSize = 0 },
		func(c *simulation.Config) { c.FoodCount = -1 },
		func(c *simulation.Config) { c.MaxHealth = 0 },
		func(c *simulation.Config) { c.FoodHealing = -1 },
		func(c *simulation.Config) { c.StarvationDamage = -1 },
		func(c *simulation.Config) { c.PredationSizeRatio = 0.5 },
		func(c *simulation.Config) { c.MutationRate = 1.1 },
		func(c *simulation.Config) { c.MutationSize = -1 },
		func(c *simulation.Config) { c.EliteCount = c.PopulationSize + 1 },
	}
	for i, change := range cases {
		c := simulation.DefaultConfig()
		change(&c)
		if _, err := simulation.New(c, 42); err == nil {
			t.Errorf("case %d accepted invalid config", i)
		}
	}
}
