package report

import (
	"aquarium/internal/simulation"
	"bytes"
	"encoding/csv"
	"errors"
	"testing"
)

type brokenWriter struct{}

func (brokenWriter) Write(p []byte) (int, error) { return 0, errors.New("write failed") }
func TestLineageCSVAndErrors(t *testing.T) {
	var b bytes.Buffer
	w, err := NewLineageCSV(&b, 42, "test")
	if err != nil {
		t.Fatal(err)
	}
	records := []simulation.LineageRecord{{ID: 1, BornIn: 1, EndedIn: 1, EndReason: "death", Fitness: 7}, {ID: 2, ParentA: 1, BornIn: 2, Elite: true, Fitness: 2}}
	if err := w.Write(records); err != nil {
		t.Fatal(err)
	}
	rows, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 || len(rows[0]) != 21 || rows[1][0] != "42" || rows[1][1] != "test" || rows[1][7] != "death" || rows[2][7] != "" || rows[2][9] != "2" {
		t.Fatal(rows)
	}
	if err := w.Write(records[:1]); err == nil {
		t.Fatal("duplicate accepted")
	}
	var output bytes.Buffer
	w, _ = NewLineageCSV(&output, 1, "x")
	before := output.String()
	if err := w.Write([]simulation.LineageRecord{records[0], records[0]}); err == nil || output.String() != before {
		t.Fatal("batch duplicate")
	}
	if _, err := NewLineageCSV(brokenWriter{}, 1, "x"); err == nil {
		t.Fatal("header error")
	}
	w.output = csv.NewWriter(brokenWriter{})
	if err := w.Write(records); err == nil {
		t.Fatal("flush error")
	}
	if _, err := NewGroupsCSV(brokenWriter{}, 1, "x"); err == nil {
		t.Fatal("group header error")
	}
	g, _ := NewGroupsCSV(&output, 1, "x")
	g.output = csv.NewWriter(brokenWriter{})
	if err := g.Write(1, "completed", records, 60); err == nil {
		t.Fatal("group flush error")
	}
}
