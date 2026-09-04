package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
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

	Color rl.Color
}

var seedColors = []rl.Color{
	rl.NewColor(240, 70, 70, OpaqueAlpha),  // coral
	rl.NewColor(245, 210, 50, OpaqueAlpha), // gold
	rl.NewColor(40, 220, 170, OpaqueAlpha), // teal
	rl.NewColor(190, 70, 230, OpaqueAlpha), // violet
}

func getRandomColor() rl.Color {
	seed := seedColors[rand.Intn(len(seedColors))]

	return rl.NewColor(
		jitterColorChannel(seed.R),
		jitterColorChannel(seed.G),
		jitterColorChannel(seed.B),
		OpaqueAlpha,
	)
}

func jitterColorChannel(channel uint8) uint8 {
	offset := rand.Intn(ColorJitter*2+1) - ColorJitter
	value := int(channel) + offset
	value = max(0, min(value, ColorChannelCount-1))

	return uint8(value)
}

func initRandomGenome() Genome {
	return Genome{
		Size:       rand.Float32()*FishSizeRange + MinFishSize,
		MaxSpeed:   rand.Float32()*FishSpeedRange + MinFishSpeed,
		Vision:     rand.Float32()*FishVisionRange + MinFishVision,
		Metabolism: rand.Float32()*MetabolismRange + MinMetabolism,

		AttractionToFood:   rand.Float32(),
		FearOfPredators:    rand.Float32(),
		AttractionToOthers: rand.Float32(),

		Color: getRandomColor(),
	}
}

// TODO: getEliteGenomes() function
