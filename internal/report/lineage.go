package report

import (
	"encoding/csv"
	"fmt"
	"io"

	"aquarium/internal/simulation"
)

type LineageCSV struct {
	output *csv.Writer
	seed   int64
	run    string
	seen   map[simulation.FishID]bool
}

func NewLineageCSV(w io.Writer, seed int64, run string) (*LineageCSV, error) {
	c := &LineageCSV{csv.NewWriter(w), seed, run, make(map[simulation.FishID]bool)}

	err := c.output.Write([]string{"seed", "run_id", "fish_id", "parent_a", "parent_b", "born_in", "ended_in", "end_reason", "elite", "fitness", "size", "max_speed", "vision", "metabolism", "food_attraction", "predator_fear", "social_attraction", "red", "green", "blue", "alpha"})
	
	if err != nil {
		return nil, err
	}
	
	c.output.Flush()
	
	return c, c.output.Error()
}

func (c *LineageCSV) Write(records []simulation.LineageRecord) error {
	batch := make(map[simulation.FishID]bool)
	
	for _, r := range records {
		if c.seen[r.ID] || batch[r.ID] {
			return fmt.Errorf("duplicate lineage ID %d", r.ID)
		}
	
		batch[r.ID] = true
	}
	
	for _, r := range records {
		g := r.Genome
		row := []string{fmt.Sprint(c.seed), c.run, fmt.Sprint(r.ID), fmt.Sprint(r.ParentA), fmt.Sprint(r.ParentB), fmt.Sprint(r.BornIn), fmt.Sprint(r.EndedIn), r.EndReason, fmt.Sprint(r.Elite), fmt.Sprint(r.Fitness), fmt.Sprint(g.Size), fmt.Sprint(g.MaxSpeed), fmt.Sprint(g.Vision), fmt.Sprint(g.Metabolism), fmt.Sprint(g.AttractionToFood), fmt.Sprint(g.FearOfPredators), fmt.Sprint(g.AttractionToOthers), fmt.Sprint(g.Color.R), fmt.Sprint(g.Color.G), fmt.Sprint(g.Color.B), fmt.Sprint(g.Color.A)}
		
		if err := c.output.Write(row); err != nil {
			return err
		}
		
		c.seen[r.ID] = true
	}
	
	c.output.Flush()
	
	return c.output.Error()
}

type GroupsCSV struct {
	output *csv.Writer
	seed   int64
	run    string
}

func NewGroupsCSV(w io.Writer, seed int64, run string) (*GroupsCSV, error) {
	
	c := &GroupsCSV{csv.NewWriter(w), seed, run}
	
	err := c.output.Write([]string{"seed", "run_id", "generation", "cohort", "threshold", "group_id", "population", "share", "mean_fitness", "red", "green", "blue"})
	if err != nil {
		return nil, err
	}
	
	c.output.Flush()
	
	return c, c.output.Error()
}

func (c *GroupsCSV) Write(generation int, cohort string, records []simulation.LineageRecord, threshold float64) error {
	sample := make([]ColorSample, len(records))
	
	for i, r := range records {
		sample[i] = ColorSample{r.ID, r.Genome.Color, r.Fitness}
	}
	
	groups, err := GroupColors(sample, threshold)
	if err != nil {
		return err
	}
	
	for _, g := range groups.Groups {
		if err := c.output.Write([]string{fmt.Sprint(c.seed), c.run, fmt.Sprint(generation), cohort, fmt.Sprint(threshold), fmt.Sprint(g.ID), fmt.Sprint(g.Population), fmt.Sprint(g.Share), fmt.Sprint(g.MeanFitness), fmt.Sprint(g.MeanColor.R), fmt.Sprint(g.MeanColor.G), fmt.Sprint(g.MeanColor.B)}); err != nil {
			return err
		}
	}
	
	c.output.Flush()
	
	return c.output.Error()
}
