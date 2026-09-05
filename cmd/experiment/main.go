package main

import (
	"aquarium/internal/report"
	"aquarium/internal/simulation"

	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func main() {
	seed := flag.Int64("seed", 42, "random seed")
	generations := flag.Int("generations", 100, "completed generations")
	population := flag.Int("population", 50, "fish per generation")
	outputPath := flag.String("output", "-", "CSV destination, or - for stdout")
	step := flag.Float64("step", 1.0/60.0, "fixed simulation step in seconds")
	defaults := simulation.DefaultConfig()
	speedEnergy := flag.Float64("speed-energy", float64(defaults.SpeedEnergyWeight), "movement hunger rate at maximum speed")
	sizeEnergy := flag.Float64("size-energy", float64(defaults.SizeEnergyWeight), "hunger rate at maximum size")
	visionEnergy := flag.Float64("vision-energy", float64(defaults.VisionEnergyWeight), "hunger rate at maximum vision")
	acceleration := flag.Float64("acceleration", float64(defaults.BaseAcceleration), "steering acceleration; zero disables limit")
	ratio := flag.Float64("predation-ratio", float64(defaults.PredationSizeRatio), "required predator/prey size ratio")
	duration := flag.Float64("duration", float64(defaults.GenerationDuration), "generation duration in seconds")
	food := flag.Int("food", defaults.FoodCount, "food count")
	flag.Parse()
	if *generations <= 0 || *step <= 0 || math.IsNaN(*step) || math.IsInf(*step, 0) || float32(*step) == 0 || math.IsInf(float64(float32(*step)), 0) {
		log.Fatal("generations and finite representable step must be positive")
	}

	config := simulation.DefaultConfig()
	config.PopulationSize = *population
	config.SpeedEnergyWeight = float32(*speedEnergy)
	config.SizeEnergyWeight = float32(*sizeEnergy)
	config.VisionEnergyWeight = float32(*visionEnergy)
	config.BaseAcceleration = float32(*acceleration)
	config.PredationSizeRatio = float32(*ratio)
	config.GenerationDuration = float32(*duration)
	config.FoodCount = *food
	sim, err := simulation.New(config, *seed)
	if err != nil {
		log.Fatal(err)
	}

	for len(sim.History()) < *generations {
		sim.Step(float32(*step))
	}

	var destination io.Writer = os.Stdout

	if *outputPath != "-" {
		outputFile, err := os.Create(*outputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()
		destination = outputFile
	}

	history := sim.History()
	if err := report.WriteCSV(destination, history); err != nil {
		log.Fatal(err)
	}

	last := history[len(history)-1]
	fmt.Fprintf(os.Stderr, "seed=%d generations=%d best_fitness=%.4f output=%s\n",
		*seed, *generations, last.BestFitness, *outputPath)
}
