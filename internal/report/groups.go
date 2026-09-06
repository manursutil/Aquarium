package report

import (
	"aquarium/internal/simulation"
	"fmt"
	"math"
	"sort"
)

const DefaultColorGroupThreshold = 60.0

type ColorSample struct {
	ID      simulation.FishID
	Color   simulation.Color
	Fitness float32
}
type ColorGroup struct {
	ID, Population     int
	Share, MeanFitness float64
	MeanColor          simulation.Color
}
type Membership struct {
	FishID  simulation.FishID
	GroupID int
}
type ColorGroups struct {
	Groups      []ColorGroup
	Memberships []Membership
}

func GroupColors(sample []ColorSample, threshold float64) (ColorGroups, error) {
	result := ColorGroups{}
	if threshold <= 0 || math.IsNaN(threshold) || math.IsInf(threshold, 0) {
		return result, fmt.Errorf("color group threshold must be finite and positive")
	}

	sample = append([]ColorSample(nil), sample...)
	sort.Slice(sample, func(i, j int) bool { return sample[i].ID < sample[j].ID })
	representatives := []simulation.Color{}
	
	totals := [][3]uint64{}
	for i, f := range sample {
		if i > 0 && f.ID == sample[i-1].ID {
			return ColorGroups{}, fmt.Errorf("duplicate fish ID %d", f.ID)
		}

		chosen, nearest := -1, threshold*threshold
		for j, c := range representatives {
			dr, dg, db := float64(f.Color.R)-float64(c.R), float64(f.Color.G)-float64(c.G), float64(f.Color.B)-float64(c.B)
			d := dr*dr + dg*dg + db*db
		
			if d <= nearest && (chosen < 0 || d < nearest) {
				chosen, nearest = j, d
			}
		}
		
		if chosen < 0 {
			chosen = len(representatives)
			representatives = append(representatives, f.Color)
			totals = append(totals, [3]uint64{})
			result.Groups = append(result.Groups, ColorGroup{ID: chosen + 1})
		}
		
		g := &result.Groups[chosen]
		g.Population++
		g.MeanFitness += float64(f.Fitness)
		
		totals[chosen][0] += uint64(f.Color.R)
		totals[chosen][1] += uint64(f.Color.G)
		totals[chosen][2] += uint64(f.Color.B)
		result.Memberships = append(result.Memberships, Membership{f.ID, chosen + 1})
	}
	
	for i := range result.Groups {
		g := &result.Groups[i]
		n := uint64(g.Population)
		
		g.Share = float64(g.Population) / float64(len(sample))
		g.MeanFitness /= float64(g.Population)
		g.MeanColor = simulation.Color{R: uint8(totals[i][0] / n), G: uint8(totals[i][1] / n), B: uint8(totals[i][2] / n), A: 255}
	}
	return result, nil
}
