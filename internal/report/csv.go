package report

import (
	"encoding/csv"
	"fmt"
	"io"

	"aquarium/internal/simulation"
)

var generationHeader = []string{
	"generation", "evaluated", "survivors", "best_fitness",
	"mean_fitness", "median_fitness", "food_eaten", "fish_eaten",
	"size_mean", "speed_mean", "vision_mean", "metabolism_mean",
	"size_min", "size_max", "speed_min", "speed_max", "vision_min", "vision_max",
}

func WriteCSV(writer io.Writer, history []simulation.GenerationStats) error {
	output := csv.NewWriter(writer)
	if err := output.Write(generationHeader); err != nil {
		return err
	}

	for _, stats := range history {
		row := []string{
			fmt.Sprint(stats.Generation),
			fmt.Sprint(stats.Evaluated),
			fmt.Sprint(stats.Survivors),
			fmt.Sprintf("%.4f", stats.BestFitness),
			fmt.Sprintf("%.4f", stats.MeanFitness),
			fmt.Sprintf("%.4f", stats.MedianFitness),
			fmt.Sprint(stats.FoodEaten),
			fmt.Sprint(stats.FishEaten),
			fmt.Sprintf("%.4f", stats.Size.Mean),
			fmt.Sprintf("%.4f", stats.MaxSpeed.Mean),
			fmt.Sprintf("%.4f", stats.Vision.Mean),
			fmt.Sprintf("%.4f", stats.Metabolism.Mean),
			fmt.Sprintf("%.4f", stats.Size.Min), fmt.Sprintf("%.4f", stats.Size.Max),
			fmt.Sprintf("%.4f", stats.MaxSpeed.Min), fmt.Sprintf("%.4f", stats.MaxSpeed.Max),
			fmt.Sprintf("%.4f", stats.Vision.Min), fmt.Sprintf("%.4f", stats.Vision.Max),
		}
		if err := output.Write(row); err != nil {
			return err
		}
	}

	output.Flush()

	return output.Error()
}
