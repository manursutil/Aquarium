package main

import (
	"aquarium/internal/simulation"
	"fmt"
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

func fishTooltipLines(fish simulation.FishSnapshot, maxHealth float32) []string {
	return []string{
		fmt.Sprintf("Health: %.1f / %.1f", fish.Health, maxHealth),
		fmt.Sprintf("Hunger: %.2f", fish.Hunger),
		fmt.Sprintf("Age: %.1fs", fish.Age),
		fmt.Sprintf("Fitness: %.1f", fish.Fitness),
		fmt.Sprintf("Food eaten: %d", fish.FoodEaten),
		fmt.Sprintf("Fish eaten: %d", fish.FishEaten),
		"",
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

func drawFishTooltip(fish simulation.FishSnapshot, mouse rl.Vector2, maxHealth float32) {
	lines := fishTooltipLines(fish, maxHealth)
	height := panelHeight(len(lines))
	position := tooltipPosition(mouse, TooltipWidth, height)

	drawInfoPanel(
		lines,
		position.X,
		position.Y,
		TooltipWidth,
		toRaylibColor(fish.Color),
	)

	drawFishSelectionRing(fish)
}

func drawPinnedFishInspector(fish simulation.FishSnapshot, maxHealth float32) {
	lines := append([]string{"Pinned fish (Esc to unpin)", ""}, fishTooltipLines(fish, maxHealth)...)
	x := max(UIBorder, float32(rl.GetScreenWidth())-TooltipWidth-UIBorder)
	drawInfoPanel(lines, x, UIBorder, TooltipWidth, toRaylibColor(fish.Color))
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

func drawHUD(s simulation.Snapshot) {
	lines := []string{
		fmt.Sprintf("Generation: %d", s.Generation),
		fmt.Sprintf("Alive: %d / %d", len(s.Fish), len(s.Fish)+s.Evaluated),
		fmt.Sprintf("Evaluated: %d", s.Evaluated),
		fmt.Sprintf("Time left: %.1fs", max(float32(0), s.GenerationDuration-s.GenerationElapsed)),
		fmt.Sprintf("Best fitness: %.1f", s.BestFitness),
		fmt.Sprintf("Food: %d", len(s.Food)),
		fmt.Sprintf("FPS: %d", rl.GetFPS()),
	}
	drawInfoPanel(lines, UIBorder, UIBorder, HUDWidth, rl.SkyBlue)
}
