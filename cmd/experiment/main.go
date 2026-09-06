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
	"path/filepath"
	"time"
)

func main() {
	lineagePath := flag.String("lineage-output", "", "optional lineage CSV file")
	groupsPath := flag.String("groups-output", "", "optional color groups CSV file")
	runID := flag.String("run-id", time.Now().UTC().Format(time.RFC3339Nano), "unique run identifier for combined exports")
	threshold := flag.Float64("color-group-threshold", report.DefaultColorGroupThreshold, "reporting-only RGB distance")
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

	if _, err := report.GroupColors(nil, *threshold); err != nil {
		log.Fatal(err)
	}
	if *lineagePath == "-" || *groupsPath == "-" {
		log.Fatal("lineage and group outputs require file paths; stdout is reserved for generation CSV")
	}
	paths := map[string]bool{}
	for _, p := range []string{*outputPath, *lineagePath, *groupsPath} {
		if p == "" || p == "-" {
			continue
		}
		absolute, err := filepath.Abs(p)
		if err != nil {
			log.Fatal(err)
		}
		if paths[absolute] {
			log.Fatal("output paths must be distinct")
		}
		paths[absolute] = true
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

	var lineage *report.LineageCSV
	var groups *report.GroupsCSV
	if *lineagePath != "" {
		f, err := os.Create(*lineagePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		lineage, err = report.NewLineageCSV(f, *seed, *runID)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *groupsPath != "" {
		f, err := os.Create(*groupsPath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		groups, err = report.NewGroupsCSV(f, *seed, *runID)
		if err != nil {
			log.Fatal(err)
		}
	}
	completed := 0
	for completed < *generations {
		sim.Step(float32(*step))
		current := len(sim.History())
		if current > completed {
			if lineage != nil || groups != nil {
				cohort := []simulation.LineageRecord{}
				for _, r := range sim.LineageRecords() {
					if r.EndedIn == current {
						cohort = append(cohort, r)
					}
				}
				if lineage != nil {
					if err := lineage.Write(cohort); err != nil {
						log.Fatal(err)
					}
				}
				if groups != nil {
					if err := groups.Write(current, "completed", cohort, *threshold); err != nil {
						log.Fatal(err)
					}
				}
			}
			completed = current
		}
	}
	active := []simulation.LineageRecord{}
	for _, r := range sim.LineageRecords() {
		if r.EndedIn == 0 {
			active = append(active, r)
		}
	}
	if lineage != nil {
		if err := lineage.Write(active); err != nil {
			log.Fatal(err)
		}
	}
	if groups != nil {
		if err := groups.Write(completed+1, "living", active, *threshold); err != nil {
			log.Fatal(err)
		}
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
