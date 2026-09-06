package simulation

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestLineageLifecycleAndRetention(t *testing.T) {
	c := DefaultConfig()
	c.PopulationSize = 3
	c.LineageGenerations = 1
	a, err := New(c, 42)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range a.fish {
		r, ok := a.Lineage(f.ID)
		if !ok || r.BornIn != 1 || r.ParentA != 0 || r.ParentB != 0 || r.Elite || r.EndReason != "" {
			t.Fatalf("founder: %+v", r)
		}
	}
	a.fish[0].Age = 7
	a.fish[0].FoodEaten = 2
	r, _ := a.Lineage(1)
	if r.Fitness != 17 {
		t.Fatal(r)
	}
	r.Genome.Size = 999
	if r, _ := a.Lineage(1); r.Genome.Size == 999 {
		t.Fatal("lookup aliases world")
	}
	a.fish[0].Alive = false
	a.removeDead()
	a.removeDead()
	if len(a.generationCandidates) != 1 || a.generationCandidates[0].FishID != 1 {
		t.Fatal(a.generationCandidates)
	}
	r, _ = a.Lineage(1)
	if r.EndReason != "death" || r.EndedIn != 1 || r.Fitness != 17 {
		t.Fatal(r)
	}
	candidates := append(append([]Candidate{}, a.generationCandidates...), makeCandidates(a.fish)...)
	seen := map[FishID]bool{}
	for _, c := range candidates {
		if seen[c.FishID] {
			t.Fatal("duplicate candidate")
		}
		seen[c.FishID] = true
	}
	if len(seen) != 3 {
		t.Fatal(seen)
	}
	a.advanceGeneration()
	if a.history[0].BestFishID != 1 {
		t.Fatal(a.history)
	}
	for _, id := range []FishID{2, 3} {
		r, _ := a.Lineage(id)
		if r.EndReason != "rollover" || r.EndedIn != 1 {
			t.Fatal(r)
		}
	}
	elite := a.fish[0]
	source, _ := a.Lineage(1)
	if elite.ID <= 3 || elite.ParentA != 1 || elite.ParentB != 0 || !elite.Elite || elite.BornIn != 2 || elite.Genome != source.Genome || elite.Age != 0 || elite.Hunger != 0 || elite.Health != c.MaxHealth || elite.FoodEaten != 0 || elite.FishEaten != 0 {
		t.Fatal(elite)
	}
	for _, f := range a.fish[1:] {
		if !seen[f.ParentA] || !seen[f.ParentB] || f.Elite {
			t.Fatal(f)
		}
	}
	records := a.LineageRecords()
	records[0].Genome.Size = 999
	if r, _ := a.Lineage(1); r.Genome.Size == 999 {
		t.Fatal("records alias world")
	}
	for range 5 {
		a.advanceGeneration()
		if len(a.LineageRecords()) != 6 {
			t.Fatal("retention bound", len(a.LineageRecords()))
		}
	}
	if _, ok := a.Lineage(1); ok {
		t.Fatal("old founder retained")
	}
	if _, ok := a.Lineage(0); ok {
		t.Fatal("zero ID exists")
	}
	if _, ok := a.Lineage(9999); ok {
		t.Fatal("unknown ID exists")
	}
}

func TestOffspringSelectionAndRandomDraws(t *testing.T) {
	c := DefaultConfig()
	c.PopulationSize = 8
	c.EliteCount = 2
	candidates := []Candidate{{FishID: 17, Genome: validTestGenome(), Fitness: 1}, {FishID: 21, Genome: validTestGenome(), Fitness: 4}}
	candidates[1].Genome.Color.R = 240
	actual, oracle := rand.New(rand.NewSource(19)), rand.New(rand.NewSource(19))
	got := evolve(actual, candidates, c)
	pick := func() Candidate {
		best := candidates[oracle.Intn(len(candidates))]
		for range 3 {
			p := candidates[oracle.Intn(len(candidates))]
			if p.Fitness > best.Fitness {
				best = p
			}
		}
		return best
	}
	for i, child := range got {
		if i < c.EliteCount {
			if child != (Offspring{Genome: candidates[1].Genome, ParentA: 21, Elite: true}) {
				t.Fatal(child)
			}
			continue
		}
		a, b := pick(), pick()
		genome := crossover(a.Genome, b.Genome)
		mutate(oracle, c, &genome)
		if child != (Offspring{Genome: genome, ParentA: a.FishID, ParentB: b.FishID}) {
			t.Fatal(child)
		}
	}
	if actual.Int63() != oracle.Int63() {
		t.Fatal("random draws changed")
	}
	single := evolve(rand.New(rand.NewSource(1)), candidates[:1], c)
	for _, child := range single[c.EliteCount:] {
		if child.ParentA != 17 || child.ParentB != 17 {
			t.Fatal("self-parent", child)
		}
	}
	if evolve(actual, nil, c) != nil {
		t.Fatal("empty candidates")
	}
}

func TestLineageReplayAndExtinction(t *testing.T) {
	c := DefaultConfig()
	c.PopulationSize = 4
	c.LineageGenerations = 2
	c.GenerationDuration = 0.1
	a, _ := New(c, 42)
	b, _ := New(c, 42)
	for range 90 {
		a.Step(1.0 / 60)
		b.Step(1.0 / 60)
		if !reflect.DeepEqual(a.LineageRecords(), b.LineageRecords()) {
			t.Fatal("lineage replay")
		}
	}
	for i := range a.fish {
		a.fish[i].Alive = false
	}
	generation := a.generation
	a.Step(0)
	for _, r := range a.LineageRecords() {
		if r.BornIn == generation && r.EndReason != "death" {
			t.Fatal(r)
		}
	}
	if a.generation != generation+1 || len(a.fish) != c.PopulationSize {
		t.Fatal("extinction rollover")
	}
	c.LineageGenerations = 0
	if _, err := New(c, 1); err == nil {
		t.Fatal("zero retention accepted")
	}
}

func TestBestFishKeepsCandidateTieOrder(t *testing.T) {
	stats := summarizeGeneration(1, 1, []Candidate{{FishID: 9, Fitness: 7}, {FishID: 2, Fitness: 7}})
	if stats.BestFishID != 9 {
		t.Fatal(stats.BestFishID)
	}
}
