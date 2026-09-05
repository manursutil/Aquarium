package simulation

import (
	"math/rand"
)

const (
	FishSizeRange     float32 = 20
	MinFishSize       float32 = 10
	FishSpeedRange    float32 = 100
	MinFishSpeed      float32 = 2
	FishVisionRange   float32 = 100
	MinFishVision     float32 = 50
	MetabolismRange   float32 = 0.1
	MinMetabolism     float32 = 0.05
	ColorChannelCount         = 256
	ColorJitter               = 15
	OpaqueAlpha               = 255
)

type Genome struct {
	Size       float32
	MaxSpeed   float32
	Vision     float32
	Metabolism float32

	AttractionToFood   float32
	FearOfPredators    float32
	AttractionToOthers float32

	Color Color
}

var seedColors = []Color{
	newColor(240, 70, 70, OpaqueAlpha),  // coral
	newColor(245, 210, 50, OpaqueAlpha), // gold
	newColor(40, 220, 170, OpaqueAlpha), // teal
	newColor(190, 70, 230, OpaqueAlpha), // violet
}

func getRandomColor(rng *rand.Rand) Color {
	seed := seedColors[rng.Intn(len(seedColors))]

	return newColor(
		jitterColorChannel(rng, seed.R),
		jitterColorChannel(rng, seed.G),
		jitterColorChannel(rng, seed.B),
		OpaqueAlpha,
	)
}

func jitterColorChannel(rng *rand.Rand, channel uint8) uint8 {
	offset := rng.Intn(ColorJitter*2+1) - ColorJitter
	value := int(channel) + offset
	value = max(0, min(value, ColorChannelCount-1))

	return uint8(value)
}

func initRandomGenome(rng *rand.Rand) Genome {
	return Genome{
		Size:       rng.Float32()*FishSizeRange + MinFishSize,
		MaxSpeed:   rng.Float32()*FishSpeedRange + MinFishSpeed,
		Vision:     rng.Float32()*FishVisionRange + MinFishVision,
		Metabolism: rng.Float32()*MetabolismRange + MinMetabolism,

		AttractionToFood:   rng.Float32(),
		FearOfPredators:    rng.Float32(),
		AttractionToOthers: rng.Float32(),

		Color: getRandomColor(rng),
	}
}
