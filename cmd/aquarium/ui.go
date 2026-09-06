package main

import (
	"aquarium/internal/simulation"
	"fmt"
	"math"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	UIBorder      float32 = 10
	UIPadding     float32 = 10
	UIFontSize    int32   = 16
	UILineHeight  float32 = 20
	HUDWidth      float32 = 220
	TooltipWidth  float32 = 250
	TooltipOffset float32 = 14
)

func panelHeight(linecount int) float32 {
	return UIPadding*2 + float32(linecount)*UILineHeight
}

func drawInfoPanel(lines []string, x float32, y float32, width float32, border rl.Color) {
	panel := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: panelHeight(len(lines)),
	}

	rl.DrawRectangleRounded(panel, 0.08, 6, rl.Fade(rl.Black, 0.85))
	rl.DrawRectangleRoundedLinesEx(panel, 0.08, 6, 1, border)

	for i, line := range lines {
		rl.DrawText(
			line,
			int32(x+UIPadding),
			int32(y+UIPadding+float32(i)*UILineHeight),
			UIFontSize,
			rl.RayWhite,
		)
	}
}

func hoveredFishIndex(fish []simulation.FishSnapshot, mouse rl.Vector2) int {
	for i := range slices.Backward(fish) {

		if rl.CheckCollisionPointCircle(mouse, toRaylibVector(fish[i].Position), fish[i].Radius) {
			return i
		}
	}

	return -1
}

func fishTooltipLines(fish simulation.FishSnapshot, maxHealth float32, generation int) []string {
	lines := []string{
		fmt.Sprintf("Fish %d | Generation %d", fish.ID, generation),
		fmt.Sprintf("Speed: %.2f", math.Hypot(float64(fish.Velocity.X), float64(fish.Velocity.Y))),
		fmt.Sprintf("Health: %.1f / %.1f", fish.Health, maxHealth),
		fmt.Sprintf("Hunger: %.2f", fish.Hunger),
		fmt.Sprintf("Age: %.1fs", fish.Age),
		fmt.Sprintf("Fitness: %.1f", fish.Fitness),
		fmt.Sprintf("Food eaten: %d", fish.FoodEaten),
		fmt.Sprintf("Fish eaten: %d", fish.FishEaten),
		"Genome",
		fmt.Sprintf("Size: %.2f", fish.Genome.Size),
		fmt.Sprintf("Max speed: %.2f", fish.Genome.MaxSpeed),
		fmt.Sprintf("Vision: %.2f", fish.Genome.Vision),
		fmt.Sprintf("Metabolism: %.3f", fish.Genome.Metabolism),
		fmt.Sprintf("Food attraction: %.2f", fish.Genome.AttractionToFood),
		fmt.Sprintf("Predator fear: %.2f", fish.Genome.FearOfPredators),
		fmt.Sprintf("Social attraction: %.2f", fish.Genome.AttractionToOthers),
		fmt.Sprintf(
			"Color: %d, %d, %d",
			fish.Genome.Color.R,
			fish.Genome.Color.G,
			fish.Genome.Color.B,
		),
	}
	details := append([]string{lines[0]}, birthLines(fish.ParentA, fish.ParentB, fish.BornIn, fish.Elite)...)
	return append(details, lines[1:]...)
}

func tooltipPosition(mouse rl.Vector2, width, height float32) rl.Vector2 {
	x := mouse.X + TooltipOffset
	y := mouse.Y + TooltipOffset
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	if x+width+UIBorder > screenWidth {
		x = mouse.X - width - TooltipOffset
	}
	if y+height+UIBorder > screenHeight {
		y = screenHeight - height - UIBorder
	}

	x = max(UIBorder, x)
	y = max(UIBorder, y)

	return rl.Vector2{X: x, Y: y}
}

func drawFishTooltip(fish simulation.FishSnapshot, mouse rl.Vector2, maxHealth float32, generation int, group string) {
	lines := append([]string{group}, fishTooltipLines(fish, maxHealth, generation)...)
	height := panelHeight(len(lines))
	width := max(TooltipWidth, measuredPanelWidth(lines))
	position := tooltipPosition(mouse, width, height)

	drawInfoPanel(
		lines,
		position.X,
		position.Y,
		width,
		toRaylibColor(fish.Color),
	)

	drawFishSelectionRing(fish)
}

func drawPinnedFishInspector(fish simulation.FishSnapshot, maxHealth float32, generation int, group string) {
	lines := append([]string{"Pinned fish (A: ancestry)", group}, fishTooltipLines(fish, maxHealth, generation)...)
	width := max(TooltipWidth, measuredPanelWidth(lines))
	x := max(UIBorder, float32(rl.GetScreenWidth())-width-UIBorder)
	drawInfoPanel(lines, x, UIBorder, width, toRaylibColor(fish.Color))
	drawFishSelectionRing(fish)
}

func drawFishSelectionRing(fish simulation.FishSnapshot) {
	rl.DrawCircleLines(
		int32(fish.Position.X),
		int32(fish.Position.Y),
		fish.Genome.Size+2,
		rl.White,
	)
}

func drawHUD(s simulation.Snapshot, view ViewState, seed int64) {
	state := "running"
	if view.Paused {
		state = "paused"
	}
	lines := []string{
		fmt.Sprintf("Seed: %d", seed),
		fmt.Sprintf("Speed: %dx | %s", view.SpeedMultiplier, state),
		fmt.Sprintf("Generation: %d", s.Generation),
		fmt.Sprintf("Alive: %d / %d", len(s.Fish), len(s.Fish)+s.Evaluated),
		fmt.Sprintf("Evaluated: %d", s.Evaluated),
		fmt.Sprintf("Time left: %.1fs", max(float32(0), s.GenerationDuration-s.GenerationElapsed)),
		fmt.Sprintf("Best fitness: %.1f", s.BestFitness),
		fmt.Sprintf("Food: %d", len(s.Food)),
		fmt.Sprintf("FPS: %d", rl.GetFPS()),
	}
	drawInfoPanel(lines, UIBorder, UIBorder, max(HUDWidth, measuredPanelWidth(lines)), rl.SkyBlue)
}

func measuredPanelWidth(lines []string) float32 {
	width := float32(0)
	for _, line := range lines {
		width = max(width, float32(rl.MeasureText(line, UIFontSize))+2*UIPadding)
	}
	return width
}

func drawControls(view ViewState) {
	onOff := func(on bool) string {
		if on {
			return "on"
		}
		return "off"
	}
	pinned := "none"
	if view.HasSelectedFish {
		pinned = fmt.Sprintf("fish %d", view.SelectedFishID)
	}
	lines := []string{
		"Space pause   1/2/3 speed   R restart   N new seed",
		"V vision   F steering   Click pin   Esc clear",
		fmt.Sprintf("Vision: %s | Steering: %s | Pinned: %s", onOff(view.ShowVision), onOff(view.ShowSteering), pinned),
		fmt.Sprintf("S hide 20x summaries: %s | A ancestry | B best", onOff(view.HideFastSummaries)),
	}
	width := measuredPanelWidth(lines)
	drawInfoPanel(lines, UIBorder, float32(rl.GetScreenHeight())-panelHeight(len(lines))-UIBorder, width, rl.SkyBlue)
}

func drawGenerationSummary(stats simulation.GenerationStats) {
	lines := []string{
		fmt.Sprintf("Generation %d complete", stats.Generation),
		fmt.Sprintf("Best: #%d | fitness %.1f", stats.BestFishID, stats.BestFitness),
		fmt.Sprintf("Mean fitness: %.1f", stats.MeanFitness),
		fmt.Sprintf("Survivors: %d / %d", stats.Survivors, stats.Evaluated),
		fmt.Sprintf("Food eaten: %d", stats.FoodEaten),
		fmt.Sprintf("Fish eaten: %d", stats.FishEaten),
	}
	drawInfoPanel(lines, UIBorder, 240, measuredPanelWidth(lines), rl.Gold)
}
