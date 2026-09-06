package main

import (
	"aquarium/internal/report"
	"aquarium/internal/simulation"
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func birthLines(parentA, parentB simulation.FishID, bornIn int, elite bool) []string {
	birth, parents := "founder", "founder"
	if elite {
		birth = "elite copy"
		parents = fmt.Sprint(parentA)
	} else if parentA != 0 || parentB != 0 {
		birth = "crossover"
		parents = fmt.Sprintf("%d, %d", parentA, parentB)
	}
	return []string{fmt.Sprintf("Born in: %d | %s", bornIn, birth), "Parents: " + parents}
}

func liveColorGroupLines(fish []simulation.FishSnapshot) map[simulation.FishID]string {
	sample := make([]report.ColorSample, len(fish))
	for i, f := range fish {
		sample[i] = report.ColorSample{ID: f.ID, Color: f.Color, Fitness: f.Fitness}
	}
	groups, _ := report.GroupColors(sample, report.DefaultColorGroupThreshold)
	lines := make(map[simulation.FishID]string, len(fish))
	for _, m := range groups.Memberships {
		g := groups.Groups[m.GroupID-1]
		lines[m.FishID] = fmt.Sprintf("Color group: %d | Share: %.0f%%", g.ID, g.Share*100)
	}
	return lines
}

func ancestryLines(id simulation.FishID, lookup func(simulation.FishID) (simulation.LineageRecord, bool)) []string {
	lines := []string{"Ancestry: 3 levels (Esc closes)"}
	visited := map[simulation.FishID]bool{}
	var visit func(simulation.FishID, int, string)
	visit = func(id simulation.FishID, depth int, branch string) {
		if id == 0 {
			return
		}
		prefix := strings.Repeat("  ", depth) + branch
		if visited[id] {
			lines = append(lines, fmt.Sprintf("%s#%d (see above)", prefix, id))
			return
		}
		visited[id] = true
		r, ok := lookup(id)
		if !ok {
			lines = append(lines, fmt.Sprintf("%s#%d: outside retained history", prefix, id))
			return
		}
		state := r.EndReason
		if state == "" {
			state = "active"
		}
		lines = append(lines, fmt.Sprintf("%s#%d | born %d | fitness %.1f | %s", prefix, id, r.BornIn, r.Fitness, state))
		birth := "founder"
		if r.Elite {
			birth = fmt.Sprintf("elite copy of #%d", r.ParentA)
		} else if r.ParentA != 0 || r.ParentB != 0 {
			birth = fmt.Sprintf("parents #%d, #%d", r.ParentA, r.ParentB)
		}
		lines = append(lines, strings.Repeat("  ", depth)+"  "+birth)
		if depth < 2 {
			visit(r.ParentA, depth+1, "A ")
			visit(r.ParentB, depth+1, "B ")
		}
	}
	visit(id, 0, "")
	return lines
}

func drawAncestry(sim *simulation.Simulation, id simulation.FishID) {
	lines := ancestryLines(id, sim.Lineage)
	width := measuredPanelWidth(lines)
	drawInfoPanel(lines, max(UIBorder, (float32(rl.GetScreenWidth())-width)/2), UIBorder, width, rl.Gold)
}
