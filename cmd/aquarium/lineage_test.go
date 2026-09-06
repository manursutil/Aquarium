package main

import (
	"aquarium/internal/simulation"
	"strings"
	"testing"
)

func TestAncestryBranchesAndBoundaries(t *testing.T) {
	records := map[simulation.FishID]simulation.LineageRecord{4: {ID: 4, ParentA: 2, ParentB: 3, BornIn: 3}, 2: {ID: 2, ParentA: 1, BornIn: 2, Elite: true}, 3: {ID: 3, ParentA: 1, ParentB: 99, BornIn: 2}, 1: {ID: 1, BornIn: 1}}
	lines := strings.Join(ancestryLines(4, func(id simulation.FishID) (simulation.LineageRecord, bool) { r, ok := records[id]; return r, ok }), "\n")
	for _, want := range []string{"parents #2, #3", "elite copy of #1", "founder", "#1 (see above)", "#99: outside retained history"} {
		if !strings.Contains(lines, want) {
			t.Fatalf("missing %q in %s", want, lines)
		}
	}
	for _, test := range []struct {
		a, b  simulation.FishID
		elite bool
		want  string
	}{{0, 0, false, "founder"}, {1, 2, false, "crossover"}, {1, 0, true, "elite copy"}} {
		if got := strings.Join(birthLines(test.a, test.b, 2, test.elite), " "); !strings.Contains(got, test.want) {
			t.Fatal(got)
		}
	}
	groups := liveColorGroupLines([]simulation.FishSnapshot{{ID: 1}, {ID: 2}})
	if groups[1] != "Color group: 1 | Share: 100%" {
		t.Fatal(groups)
	}
}
