package integration_test

import (
	"aquarium/internal/report"
	"aquarium/internal/simulation"
	"bytes"
	"encoding/csv"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLineageExportAcrossEvictionAndDisplayThresholds(t *testing.T) {
	c := simulation.DefaultConfig()
	c.PopulationSize = 6
	c.LineageGenerations = 1
	c.GenerationDuration = 0.05
	a, _ := simulation.New(c, 42)
	b, _ := simulation.New(c, 42)
	var output bytes.Buffer
	writer, err := report.NewLineageCSV(&output, 42, "replay")
	if err != nil {
		t.Fatal(err)
	}
	complete := 0
	for step := 0; complete < 5 && step < 100; step++ {
		a.Step(1.0 / 60)
		b.Step(1.0 / 60)
		for _, pair := range []struct {
			sim       *simulation.Simulation
			threshold float64
		}{{a, 1}, {b, 300}} {
			sample := []report.ColorSample{}
			for _, f := range pair.sim.Snapshot().Fish {
				sample = append(sample, report.ColorSample{ID: f.ID, Color: f.Color, Fitness: f.Fitness})
			}
			if _, err := report.GroupColors(sample, pair.threshold); err != nil {
				t.Fatal(err)
			}
		}
		if !reflect.DeepEqual(a.Snapshot(), b.Snapshot()) || !reflect.DeepEqual(a.History(), b.History()) || !reflect.DeepEqual(a.LineageRecords(), b.LineageRecords()) {
			t.Fatalf("display affected step %d", step)
		}
		generation := a.Snapshot().Generation - 1
		if generation > complete {
			records := []simulation.LineageRecord{}
			for _, r := range a.LineageRecords() {
				if r.EndedIn == generation {
					records = append(records, r)
				}
			}
			if len(records) != 6 {
				t.Fatal("missing cohort")
			}
			if err := writer.Write(records); err != nil {
				t.Fatal(err)
			}
			complete = generation
		}
	}
	if complete != 5 {
		t.Fatal("no rollover")
	}
	active := []simulation.LineageRecord{}
	for _, r := range a.LineageRecords() {
		if r.EndedIn == 0 {
			active = append(active, r)
		}
	}
	if err := writer.Write(active); err != nil {
		t.Fatal(err)
	}
	rows, err := csv.NewReader(&output).ReadAll()
	if err != nil || len(rows) != 37 {
		t.Fatal(len(rows), err)
	}
	ids := map[string]bool{}
	for _, r := range rows[1:] {
		if ids[r[2]] {
			t.Fatal("duplicate")
		}
		ids[r[2]] = true
		for _, p := range r[3:5] {
			if p != "0" && !ids[p] {
				t.Fatal("missing ancestor", p)
			}
		}
	}
	if _, ok := a.Lineage(1); ok {
		t.Fatal("export test did not cross eviction")
	}
}

func TestExperimentLineageFilesReplay(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "experiment")
	build := exec.Command("go", "build", "-o", binary, "./cmd/experiment")
	build.Dir = "../.."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("%s: %v", out, err)
	}
	run := func(suffix string) [][]byte {
		paths := []string{filepath.Join(dir, "stats"+suffix), filepath.Join(dir, "lineage"+suffix), filepath.Join(dir, "groups"+suffix)}
		cmd := exec.Command(binary, "-seed", "42", "-run-id", "fixture", "-population", "6", "-generations", "12", "-duration", "0.05", "-output", paths[0], "-lineage-output", paths[1], "-groups-output", paths[2])
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s: %v", out, err)
		}
		data := [][]byte{}
		for _, p := range paths {
			b, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			data = append(data, b)
		}
		return data
	}
	a, b := run("a"), run("b")
	if !reflect.DeepEqual(a, b) {
		t.Fatal("CSV replay differs")
	}
	rows, err := csv.NewReader(bytes.NewReader(a[1])).ReadAll()
	if err != nil || len(rows) != 79 {
		t.Fatal("lineage rows", len(rows), err)
	}
	groups, err := csv.NewReader(bytes.NewReader(a[2])).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if groups[1][3] != "completed" || groups[len(groups)-1][3] != "living" {
		t.Fatal("cohort labels", groups)
	}
	cmd := exec.Command(binary, "-output", filepath.Join(dir, "same"), "-lineage-output", filepath.Join(dir, "same"))
	if cmd.Run() == nil {
		t.Fatal("duplicate output path accepted")
	}
}
