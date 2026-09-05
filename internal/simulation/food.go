package simulation

import (
	"math/rand"
)

type Food struct {
	Position   Vector2
	HealAmount float32
}

func (f *Food) reposition(rng *rand.Rand, config Config) {
	f.Position = Vector2{
		X: rng.Float32() * config.WorldWidth,
		Y: rng.Float32() * config.WorldHeight,
	}
}
