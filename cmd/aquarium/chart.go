package main

import (
	"aquarium/internal/simulation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type chartPoint struct {
	X float32
	Y float32
}

func fitnessPoints(history []simulation.GenerationStats, width, height float32) []chartPoint {
	if len(history) == 0 {
		return nil
	}

	minValue := min(history[0].MeanFitness, history[0].MedianFitness)
	maxValue := max(history[0].BestFitness, history[0].MedianFitness)
	for _, stats := range history[1:] {
		minValue = min(minValue, stats.MeanFitness, stats.MedianFitness)
		maxValue = max(maxValue, stats.BestFitness, stats.MedianFitness)
	}
	if maxValue == minValue {
		maxValue = minValue + 1
	}

	points := make([]chartPoint, len(history))
	for i, stats := range history {
		x := float32(i) / float32(max(1, len(history)-1)) * width
		normalized := (stats.MeanFitness - minValue) / (maxValue - minValue)
		points[i] = chartPoint{X: x, Y: height - normalized*height}
	}
	return points
}

func visibleHistory(history []simulation.GenerationStats, limit int) []simulation.GenerationStats {
	if len(history) <= limit {
		return history
	}
	return history[len(history)-limit:]
}

func drawFitnessChart(history []simulation.GenerationStats, bounds rl.Rectangle) {
	if len(history) == 0 || bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	history = visibleHistory(history, 60)

	background := rl.Fade(rl.Black, 0.78)
	border := rl.Fade(rl.SkyBlue, 0.8)
	rl.DrawRectangleRec(bounds, background)
	rl.DrawRectangleLinesEx(bounds, 1, border)

	const (
		padding    float32 = 8
		headerSize int32   = 12
		fontSize   int32   = 10
	)

	rl.DrawText("Fitness", int32(bounds.X+padding), int32(bounds.Y+padding), headerSize, rl.RayWhite)

	series := []struct {
		label string
		color rl.Color
		value func(simulation.GenerationStats) float32
	}{
		{label: "best", color: rl.Yellow, value: func(s simulation.GenerationStats) float32 { return s.BestFitness }},
		{label: "mean", color: rl.SkyBlue, value: func(s simulation.GenerationStats) float32 { return s.MeanFitness }},
		{label: "median", color: rl.LightGray, value: func(s simulation.GenerationStats) float32 { return s.MedianFitness }},
	}

	legendX := bounds.X + 72
	for _, item := range series {
		rl.DrawLineEx(
			rl.Vector2{X: legendX, Y: bounds.Y + padding + 6},
			rl.Vector2{X: legendX + 12, Y: bounds.Y + padding + 6},
			2,
			item.color,
		)
		rl.DrawText(item.label, int32(legendX+16), int32(bounds.Y+padding+1), fontSize, item.color)
		legendX += float32(38 + len(item.label)*5)
	}

	plot := rl.Rectangle{
		X:      bounds.X + padding,
		Y:      bounds.Y + 26,
		Width:  bounds.Width - padding*2,
		Height: bounds.Height - 34,
	}
	if plot.Width <= 0 || plot.Height <= 0 {
		return
	}

	minValue := min(history[0].MeanFitness, history[0].MedianFitness)
	maxValue := max(history[0].BestFitness, history[0].MedianFitness)
	for _, stats := range history[1:] {
		minValue = min(minValue, stats.MeanFitness, stats.MedianFitness)
		maxValue = max(maxValue, stats.BestFitness, stats.MedianFitness)
	}
	if maxValue == minValue {
		maxValue = minValue + 1
	}

	pointsFor := func(value func(simulation.GenerationStats) float32) []chartPoint {
		points := make([]chartPoint, len(history))
		for i, stats := range history {
			x := float32(i) / float32(max(1, len(history)-1)) * plot.Width
			normalized := (value(stats) - minValue) / (maxValue - minValue)
			points[i] = chartPoint{X: x, Y: plot.Height - normalized*plot.Height}
		}
		return points
	}

	for i := 1; i < 4; i++ {
		y := plot.Y + float32(i)*plot.Height/4
		rl.DrawLineEx(
			rl.Vector2{X: plot.X, Y: y},
			rl.Vector2{X: plot.X + plot.Width, Y: y},
			1,
			rl.Fade(rl.Gray, 0.45),
		)
	}

	for _, item := range series {
		points := fitnessPoints(history, plot.Width, plot.Height)
		if item.label != "mean" {
			points = pointsFor(item.value)
		}

		for i := 1; i < len(points); i++ {
			rl.DrawLineEx(
				rl.Vector2{X: plot.X + points[i-1].X, Y: plot.Y + points[i-1].Y},
				rl.Vector2{X: plot.X + points[i].X, Y: plot.Y + points[i].Y},
				2,
				item.color,
			)
		}
	}
}
