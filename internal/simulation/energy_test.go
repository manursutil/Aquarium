package simulation

import (
	"math"
	"testing"
)

func near(t *testing.T, got, want float32) {
	t.Helper()
	if math.IsNaN(float64(got)) || math.Abs(float64(got-want)) > 1e-5 {
		t.Fatalf("got %g, want %g", got, want)
	}
}

func TestNormalizedTraits(t *testing.T) {
	for _, tc := range []struct{ value, want float32 }{{-1, 0}, {0, 0}, {0.5, 0.5}, {1, 1}, {2, 1}} {
		near(t, normalizedSize(MinFishSize+tc.value*FishSizeRange), tc.want)
		near(t, normalizedVision(MinFishVision+tc.value*FishVisionRange), tc.want)
	}
	near(t, normalizedTrait(10, 0, 0), 0)
	near(t, normalizedTrait(10, 0, -1), 0)
	for _, tc := range []struct{ speed, want float32 }{{0, 0}, {maximumPossibleSpeed / 2, 0.5}, {maximumPossibleSpeed, 1}, {maximumPossibleSpeed * 2, 1}} {
		near(t, normalizedSpeed(Vector2{tc.speed, 0}), tc.want)
	}
}

func TestEnergyBreakdown(t *testing.T) {
	c := DefaultConfig()
	f := Fish{Genome: Genome{Size: MinFishSize, Vision: MinFishVision, Metabolism: 0.1, MaxSpeed: maximumPossibleSpeed}}
	if got := energyCosts(f, c); got != (EnergyCosts{}) {
		t.Fatalf("minimum stationary cost: %+v", got)
	}
	near(t, hungerRate(f, c), 0.1)
	f.Genome.Size += FishSizeRange
	f.Genome.Vision += FishVisionRange
	f.Velocity = Vector2{maximumPossibleSpeed / 2, 0}
	half := energyCosts(f, c)
	near(t, half.Movement, 0.05)
	near(t, half.Size, 0.1)
	near(t, half.Vision, 0.1)
	f.Velocity.X *= 2
	full := energyCosts(f, c)
	near(t, full.Movement, 4*half.Movement)
	near(t, full.Size, half.Size)
	near(t, full.Vision, half.Vision)
	near(t, full.Total, full.Movement+full.Size+full.Vision)
	near(t, hungerRate(f, c), 0.5)
	c.SpeedEnergyWeight = 0
	c.SizeEnergyWeight = 0
	c.VisionEnergyWeight = 0
	near(t, hungerRate(f, c), 0.1)
}

func TestEnergyHungerTimeAndBounds(t *testing.T) {
	c := DefaultConfig()
	f := Fish{Genome: Genome{Size: 30, Vision: 150, Metabolism: 0.1}, Health: c.MaxHealth}
	g := f
	f.handleHunger(1, c)
	g.handleHunger(0.5, c)
	g.handleHunger(0.5, c)
	near(t, f.Hunger, 0.3)
	near(t, g.Hunger, f.Hunger)
	if f.Genome != g.Genome {
		t.Fatal("hunger changed genome")
	}
	f.handleHunger(100, c)
	near(t, f.Hunger, MaxHunger)
	near(t, f.Health, 0)
}

func TestEnergyConfiguration(t *testing.T) {
	for _, field := range []string{"speed", "size", "vision", "acceleration"} {
		for _, value := range []float32{-1, float32(math.NaN()), float32(math.Inf(1)), float32(math.Inf(-1)), 0, 1} {
			c := DefaultConfig()
			switch field {
			case "speed":
				c.SpeedEnergyWeight = value
			case "size":
				c.SizeEnergyWeight = value
			case "vision":
				c.VisionEnergyWeight = value
			case "acceleration":
				c.BaseAcceleration = value
			}
			err := c.Validate()
			valid := value >= 0 && !math.IsInf(float64(value), 0)
			if (err == nil) != valid {
				t.Errorf("%s=%g: %v", field, value, err)
			}
		}
	}
}

func TestEnergyConfigurationRejectsOverflow(t *testing.T) {
	c := DefaultConfig()
	c.SpeedEnergyWeight = math.MaxFloat32
	c.SizeEnergyWeight = math.MaxFloat32
	if c.Validate() == nil {
		t.Fatal("accepted overflowing hunger rate")
	}
}
