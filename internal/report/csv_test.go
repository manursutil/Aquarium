package report

import (
	"bytes"
	"testing"

	"aquarium/internal/simulation"
)

func TestWriteCSV(t *testing.T) {
	history := []simulation.GenerationStats{{
		Generation: 1, Evaluated: 2, Survivors: 1,
		BestFitness: 8, MeanFitness: 6, MedianFitness: 6,
		Size: simulation.TraitStats{Mean: 12},
	}}

	var buffer bytes.Buffer
	if err := WriteCSV(&buffer, history); err != nil {
		t.Fatal(err)
	}

	want := "generation,evaluated,survivors,best_fitness,mean_fitness,median_fitness,food_eaten,fish_eaten,size_mean,speed_mean,vision_mean,metabolism_mean\n" +
		"1,2,1,8.0000,6.0000,6.0000,0,0,12.0000,0.0000,0.0000,0.0000\n"
	if buffer.String() != want {
		t.Fatalf("CSV = %q, want %q", buffer.String(), want)
	}
}
