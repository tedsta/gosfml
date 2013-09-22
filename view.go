package sf

import "math"

type View struct {
	center              Vector2   // Center of the view, in scene coordinates
	size                Vector2   // Size of the view, in scene coordinates
	rot                 float32   // Angle of rotation of the view rect, in degrees
	viewport            Rect      // Viewport rectangle, expressed as a factor of the render-target's size
	transform           Transform // Precomputed projection transform corresponding to the view
	invTransform        Transform // Precomputed inverse projection transform corresponding to the view
	transformUpdated    bool      // Internal state telling if the transform needs to be updated
	invTransformUpdated bool      // Internal state telling if the inverse transform needs to be updated
}

func NewView() *View {
	return &View{Vector2{}, Vector2{}, 0, Rect{0, 0, 1, 1}, IdentityTransform(),
		IdentityTransform(), false, false}
}

func (v *View) SetCenterXY(x, y float32) {
	v.center.X = x
	v.center.Y = y

	v.transformUpdated = false
	v.invTransformUpdated = false
}

func (v *View) SetCenter(center Vector2) {
	v.SetCenterXY(center.X, center.Y)
}

func (v *View) SetSizeXY(w, h float32) {
	v.size.X = w
	v.size.Y = h

	v.transformUpdated = false
	v.invTransformUpdated = false
}

func (v *View) SetSize(size Vector2) {
	v.SetSizeXY(size.X, size.Y)
}

func (v *View) SetRotation(angle float32) {
	v.rot = angle
	for v.rot >= 360 {
		v.rot -= 360
	}
	for v.rot < 0 {
		v.rot += 360
	}

	v.transformUpdated = false
	v.invTransformUpdated = false
}

func (v *View) SetViewport(viewport Rect) {
	v.viewport = viewport
}

func (v *View) Reset(rect Rect) {
	v.center.X = rect.Left + rect.W/2
	v.center.Y = rect.Top + rect.H/2
	v.size.X = rect.W
	v.size.Y = rect.H
	v.rot = 0

	v.transformUpdated = false
	v.invTransformUpdated = false
}

func (v *View) Center() Vector2 {
	return v.center
}

func (v *View) Size() Vector2 {
	return v.size
}

func (v *View) Rotation() float32 {
	return v.rot
}

func (v *View) Viewport() Rect {
	return v.viewport
}

func (v *View) MoveXY(offsetX, offsetY float32) {
	v.SetCenterXY(v.center.X+offsetX, v.center.Y+offsetY)
}

func (v *View) Move(offset Vector2) {
	v.SetCenter(v.center.Add(offset))
}

func (v *View) Rotate(angle float32) {
	v.SetRotation(v.rot + angle)
}

func (v *View) Zoom(factor float32) {
	v.SetSizeXY(v.size.X*factor, v.size.Y*factor)
}

func (v *View) Transform() Transform {
	// Recompute the matrix if needed
	if !v.transformUpdated {
		// Rotation components
		angle := v.rot * math.Pi / 180
		cosine := float32(math.Cos(float64(angle)))
		sine := float32(math.Sin(float64(angle)))
		tx := -v.center.X*cosine - v.center.Y*sine + v.center.X
		ty := v.center.X*sine - v.center.Y*cosine + v.center.Y

		// Projection components
		a := 2 / v.size.X
		b := -2 / v.size.Y
		c := -a * v.center.X
		d := -b * v.center.Y

		// Rebuild the projection matrix
		v.transform = NewTransformFrom3x3(a*cosine, a*sine, a*tx+c,
			-b*sine, b*cosine, b*ty+d,
			0, 0, 1)
		v.transformUpdated = true
	}

	return v.transform
}

func (v *View) InverseTransform() Transform {
	// Recompute the matrix if needed
	if !v.invTransformUpdated {
		v.invTransform = v.Transform().Inverse()
		v.invTransformUpdated = true
	}

	return v.invTransform
}
