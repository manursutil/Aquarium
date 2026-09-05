package simulation

const (
	maximumPossibleSpeed = MinFishSpeed + FishSpeedRange
)

type EnergyCosts struct {
	Movement float32
	Size     float32
	Vision   float32
	Total    float32
}

func normalizedTrait(value float32, minimum float32, span float32) float32 {
	if span <= 0 {
		return 0
	}

	return max(0, min(1, (value-minimum)/span))
}

func normalizedSize(size float32) float32 {
	return normalizedTrait(size, MinFishSize, FishSizeRange)
}

func normalizedVision(vision float32) float32 {
	return normalizedTrait(vision, MinFishVision, FishVisionRange)
}

func normalizedSpeed(velocity Vector2) float32 {
	return max(0, min(1, length(velocity)/maximumPossibleSpeed))
}

func energyCosts(f Fish, config Config) EnergyCosts {
	speed := normalizedSpeed(f.Velocity)

	costs := EnergyCosts{
		Movement: config.SpeedEnergyWeight * speed * speed,
		Size:     config.SizeEnergyWeight * normalizedSize(f.Genome.Size),
		Vision:   config.VisionEnergyWeight * normalizedVision(f.Genome.Vision),
	}
	costs.Total = costs.Movement + costs.Size + costs.Vision

	return costs
}

func hungerRate(f Fish, config Config) float32 {
	return f.Genome.Metabolism + energyCosts(f, config).Total
}
