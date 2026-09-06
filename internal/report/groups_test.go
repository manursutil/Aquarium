package report

import (
	"aquarium/internal/simulation"
	"math"
	"reflect"
	"testing"
)

func TestColorGroupsRules(t *testing.T) {
	sample := []ColorSample{{1, simulation.Color{R: 0, A: 255}, 2}, {2, simulation.Color{R: 20, A: 255}, 6}, {3, simulation.Color{R: 10, A: 255}, 4}}
	g, err := GroupColors(sample, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Groups) != 2 || g.Memberships[2].GroupID != 1 || g.Groups[0].Population != 2 || g.Groups[0].MeanColor.R != 5 || g.Groups[0].MeanFitness != 3 {
		t.Fatal(g)
	}
	shuffled := []ColorSample{sample[2], sample[0], sample[1]}
	other, _ := GroupColors(shuffled, 10)
	if !reflect.DeepEqual(g, other) || shuffled[0].ID != 3 {
		t.Fatal("order or input mutation")
	}
	total, share := 0, 0.0
	for _, group := range g.Groups {
		total += group.Population
		share += group.Share
	}
	if total != 3 || math.Abs(share-1) > 1e-12 {
		t.Fatal(total, share)
	}
	same := []ColorSample{{1, simulation.Color{R: 4}, 1}, {2, simulation.Color{R: 4}, 3}}
	g, _ = GroupColors(same, 10)
	if len(g.Groups) != 1 || g.Groups[0].MeanFitness != 2 || g.Groups[0].Share != 1 {
		t.Fatal(g)
	}
	empty, err := GroupColors(nil, 10)
	if err != nil || len(empty.Groups) != 0 || len(empty.Memberships) != 0 {
		t.Fatal(empty, err)
	}
	for _, threshold := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		if _, err := GroupColors(nil, threshold); err == nil {
			t.Fatal("invalid threshold")
		}
	}
	if _, err := GroupColors([]ColorSample{sample[0], sample[0]}, 10); err == nil {
		t.Fatal("duplicate ID")
	}
	// Fixed representative at 0: the third color must start a second group,
	// even though it is within threshold of the second member and updated mean.
	fixed := []ColorSample{{ID: 1}, {ID: 2, Color: simulation.Color{R: 10}}, {ID: 3, Color: simulation.Color{R: 14}}}
	g, _ = GroupColors(fixed, 10)
	if len(g.Groups) != 2 {
		t.Fatal("representative moved", g)
	}
}
