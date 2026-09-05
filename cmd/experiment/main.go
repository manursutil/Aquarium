package main

import (
	"aquarium/internal/report"
	"aquarium/internal/simulation"

	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	seed := flag.Int64("seed", 42, "random seed")
	generations := flag.Int("generations", 100, "completed generations")
	population := flag.Int("population", 50, "fish per generation")
	outputPath := flag.String("output", "-", "CSV destination, or - for stdout")
	step := flag.Float64("step", 1.0/60.0, "fixed simulation step in seconds")
	flag.Parse()

	config := simulation.DefaultConfig()
	config.PopulationSize = *population
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
