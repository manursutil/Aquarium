package simulation

import "math"

type Vector2 struct {
	X float32
	Y float32
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func newColor(r, g, b, a uint8) Color {
	return Color{r, g, b, a}
}

func add(a, b Vector2) Vector2 {
	return Vector2{a.X + b.X, a.Y + b.Y}
}

func subtract(a, b Vector2) Vector2 {
	return Vector2{a.X - b.X, a.Y - b.Y}
}

func scale(v Vector2, f float32) Vector2 {
	return Vector2{v.X * f, v.Y * f}
}

func length(v Vector2) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func distance(a, b Vector2) float32 {
	return length(subtract(a, b))
}

func normalize(v Vector2) Vector2 {
	if n := length(v); n > 0 {
		return scale(v, 1/n)
	}
	return v
}

func clampMagnitude(v Vector2, low, high float32) Vector2 {
	n := length(v)
	if n > 0 {
		if n < low {
			return scale(v, low/n)
		}
		if n > high {
			return scale(v, high/n)
		}
	}
	return v
}
