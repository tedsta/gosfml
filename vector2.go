package sf

import (
	"math"
)

type Vector2 struct {
	X float32
	Y float32
}

func (v Vector2) Add(r Vector2) Vector2 {
	return Vector2{v.X + r.X, v.Y + r.Y}
}

func (v Vector2) Sub(r Vector2) Vector2 {
	return Vector2{v.X - r.X, v.Y - r.Y}
}

func (v Vector2) Mult(s float32) Vector2 {
	return Vector2{v.X * s, v.Y * s}
}

func (v Vector2) Div(s float32) Vector2 {
	return Vector2{v.X / s, v.Y / s}
}

func (v Vector2) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}
