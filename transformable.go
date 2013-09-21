package sf

import (
	"math"
)

type Transformable struct {
	origin                 Vector2
	pos                    Vector2
	rot                    float32
	scale                  Vector2
	transform              Transform
	transformNeedUpdate    bool
	invTransform           Transform
	invTransformNeedUpdate bool
}

func NewTransformable() *Transformable {
	return &Transformable{
		Vector2{0, 0},
		Vector2{0, 0},
		0,
		Vector2{1, 1},
		IdentityTransform(),
		true,
		IdentityTransform(),
		true,
	}
}

func (t *Transformable) SetPositionXY(x, y float32) {
	t.pos.X = x
	t.pos.Y = y
	t.transformNeedUpdate = true
	t.invTransformNeedUpdate = true
}

func (t *Transformable) SetPosition(pos Vector2) {
	t.SetPositionXY(pos.X, pos.Y)
}

func (t *Transformable) SetRotation(angle float32) {
	t.rot = angle
	for t.rot >= 360 {
		t.rot -= 360
	}
	for t.rot < 0 {
		t.rot += 360
	}

	t.transformNeedUpdate = true
	t.invTransformNeedUpdate = true
}

func (t *Transformable) SetScaleXY(factorX, factorY float32) {
	t.scale.X = factorX
	t.scale.Y = factorY
	t.transformNeedUpdate = true
	t.invTransformNeedUpdate = true
}

func (t *Transformable) SetScale(factors Vector2) {
	t.SetScaleXY(factors.X, factors.Y)
}

func (t *Transformable) SetOriginXY(x, y float32) {
	t.origin.X = x
	t.origin.Y = y
	t.transformNeedUpdate = true
	t.invTransformNeedUpdate = true
}

func (t *Transformable) SetOrigin(origin Vector2) {
	t.SetOriginXY(origin.X, origin.Y)
}

func (t *Transformable) Position() Vector2 {
	return t.pos
}

func (t *Transformable) Rotation() float32 {
	return t.rot
}

func (t *Transformable) Scale() Vector2 {
	return t.scale
}

func (t *Transformable) Origin() Vector2 {
	return t.origin
}

func (t *Transformable) MoveXY(offsetX, offsetY float32) {
	t.SetPositionXY(t.pos.X+offsetX, t.pos.Y+offsetY)
}

func (t *Transformable) Move(offset Vector2) {
	t.SetPositionXY(t.pos.X+offset.X, t.pos.Y+offset.Y)
}

func (t *Transformable) Rotate(angle float32) {
	t.SetRotation(t.rot + angle)
}

func (t *Transformable) ScaleXY(factorX, factorY float32) {
	t.SetScaleXY(t.scale.X*factorX, t.scale.Y*factorY)
}

func (t *Transformable) ScaleBy(factor Vector2) {
	t.SetScaleXY(t.scale.X*factor.X, t.scale.Y*factor.Y)
}

func (t *Transformable) Transform() Transform {
	// Recompute the combined transform if needed
	if t.transformNeedUpdate {
		angle := -t.rot * math.Pi / 180
		cosine := float32(math.Cos(float64(angle)))
		sine := float32(math.Sin(float64(angle)))
		sxc := t.scale.X * cosine
		syc := t.scale.Y * cosine
		sxs := t.scale.X * sine
		sys := t.scale.Y * sine
		tx := -t.origin.X*sxc - t.origin.Y*sys + t.pos.X
		ty := t.origin.X*sxs - t.origin.Y*syc + t.pos.Y

		t.transform = NewTransformFrom3x3(sxc, sys, tx,
			-sxs, syc, ty,
			0, 0, 1)
		t.transformNeedUpdate = false
	}

	return t.transform
}

func (t *Transformable) InverseTransform() Transform {
	// Recompute the inverse transform if needed
	if t.invTransformNeedUpdate {
		t.invTransform = t.Transform().Inverse()
		t.invTransformNeedUpdate = false
	}

	return t.invTransform
}
