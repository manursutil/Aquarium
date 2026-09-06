package simulation

import (
	"fmt"
	"testing"
)

// Measure initialization plus one simulated second, with the same seed each time.
// Population can decline through predation; rendering and snapshots are excluded.
func BenchmarkCohort60Steps(b *testing.B) {
	for _, population := range []int{50, 100, 250, 500} {
		b.Run(fmt.Sprintf("population_%d", population), func(b *testing.B) {
			config := DefaultConfig()
			config.PopulationSize = population
			b.ReportAllocs()
			for b.Loop() {
				world, err := New(config, 42)
				if err != nil {
					b.Fatal(err)
				}
				for range 60 {
					world.Step(1.0 / 60.0)
				}
			}
		})
	}
}
