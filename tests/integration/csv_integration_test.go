package integration_test

import (
	"bytes"
	"testing"

	"aquarium/internal/report"
	"aquarium/internal/simulation"
)

func runCSV(t *testing.T, seed int64) []byte {
	t.Helper()
	config := simulation.DefaultConfig()
	config.GenerationDuration = 1
	sim, err := simulation.New(config, seed)
	if err != nil {
		t.Fatal(err)
	}
	for len(sim.History()) < 10 {
		sim.Step(0.05)
	}

	var output bytes.Buffer
	if err := report.WriteCSV(&output, sim.History()); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func TestSameSeedProducesSameCSV(t *testing.T) {
	first := runCSV(t, 42)
	second := runCSV(t, 42)
	if !bytes.Equal(first, second) {
		t.Fatal("same seed produced different CSV output")
	}
}

func TestDifferentSeedProducesDifferentCSV(t *testing.T) {
	first := runCSV(t, 42)
	second := runCSV(t, 43)
	if bytes.Equal(first, second) {
		t.Fatal("different seeds produced identical CSV output")
	}
}
