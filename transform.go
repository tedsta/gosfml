package sf

import "math"

type Transform struct {
	Matrix [16]float32
}

// Identity matrix
var identityMatrix = [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}

func IdentityTransform() Transform {
	return Transform{identityMatrix}
}

func NewTransformFrom3x3(a00, a01, a02, a10, a11, a12, a20, a21, a22 float32) Transform {
	var matrix [16]float32
	matrix[0] = a00
	matrix[1] = a10
	matrix[2] = 0
	matrix[3] = a20
	matrix[4] = a01
	matrix[5] = a11
	matrix[6] = 0
	matrix[7] = a21
	matrix[8] = 0
	matrix[9] = 0
	matrix[10] = 1
	matrix[11] = 0
	matrix[12] = a02
	matrix[13] = a12
	matrix[14] = 0
	matrix[15] = a22
	return Transform{matrix}
}

func (t Transform) Inverse() Transform {
	// Compute the determinant
	det := t.Matrix[0]*(t.Matrix[15]*t.Matrix[5]-t.Matrix[7]*t.Matrix[13]) -
		t.Matrix[1]*(t.Matrix[15]*t.Matrix[4]-t.Matrix[7]*t.Matrix[12]) +
		t.Matrix[3]*(t.Matrix[13]*t.Matrix[4]-t.Matrix[5]*t.Matrix[12])

	// Compute the inverse if the determinant is not zero
	// (don't use an epsilon because the determinant may *really* be tiny)
	if det != 0 {
		return NewTransformFrom3x3((t.Matrix[15]*t.Matrix[5]-t.Matrix[7]*t.Matrix[13])/det,
			-(t.Matrix[15]*t.Matrix[4]-t.Matrix[7]*t.Matrix[12])/det,
			(t.Matrix[13]*t.Matrix[4]-t.Matrix[5]*t.Matrix[12])/det,
			-(t.Matrix[15]*t.Matrix[1]-t.Matrix[3]*t.Matrix[13])/det,
			(t.Matrix[15]*t.Matrix[0]-t.Matrix[3]*t.Matrix[12])/det,
			-(t.Matrix[13]*t.Matrix[0]-t.Matrix[1]*t.Matrix[12])/det,
			(t.Matrix[7]*t.Matrix[1]-t.Matrix[3]*t.Matrix[5])/det,
			-(t.Matrix[7]*t.Matrix[0]-t.Matrix[3]*t.Matrix[4])/det,
			(t.Matrix[5]*t.Matrix[0]-t.Matrix[1]*t.Matrix[4])/det)
	} else {
		return IdentityTransform()
	}
}

func (t *Transform) TransformPointXY(x, y float32) Vector2 {
	return Vector2{t.Matrix[0]*x + t.Matrix[4]*y + t.Matrix[12],
		t.Matrix[1]*x + t.Matrix[5]*y + t.Matrix[13]}
}

func (t *Transform) TransformPoint(point Vector2) Vector2 {
	return t.TransformPointXY(point.X, point.Y)
}

func (t *Transform) TransformRect(rect Rect) Rect {
	// Transform the 4 corners of the rectangle
	points := [4]Vector2{
		t.TransformPointXY(rect.Left, rect.Top),
		t.TransformPointXY(rect.Left, rect.Top+rect.H),
		t.TransformPointXY(rect.Left+rect.W, rect.Top),
		t.TransformPointXY(rect.Left+rect.W, rect.Top+rect.H),
	}

	// Compute the bounding rectangle of the transformed points
	left := points[0].X
	top := points[0].Y
	right := points[0].X
	bottom := points[0].Y
	for i := 1; i < 4; i++ {
		if points[i].X < left {
			left = points[i].X
		} else if points[i].X > right {
			right = points[i].X
		}
		if points[i].Y < top {
			top = points[i].Y
		} else if points[i].Y > bottom {
			bottom = points[i].Y
		}
	}

	return Rect{left, top, right - left, bottom - top}
}

func (t *Transform) Combine(t2 Transform) {
	*t = NewTransformFrom3x3(
		t.Matrix[0]*t2.Matrix[0]+t.Matrix[4]*t2.Matrix[1]+t.Matrix[12]*t2.Matrix[3],
		t.Matrix[0]*t2.Matrix[4]+t.Matrix[4]*t2.Matrix[5]+t.Matrix[12]*t2.Matrix[7],
		t.Matrix[0]*t2.Matrix[12]+t.Matrix[4]*t2.Matrix[13]+t.Matrix[12]*t2.Matrix[15],
		t.Matrix[1]*t2.Matrix[0]+t.Matrix[5]*t2.Matrix[1]+t.Matrix[13]*t2.Matrix[3],
		t.Matrix[1]*t2.Matrix[4]+t.Matrix[5]*t2.Matrix[5]+t.Matrix[13]*t2.Matrix[7],
		t.Matrix[1]*t2.Matrix[12]+t.Matrix[5]*t2.Matrix[13]+t.Matrix[13]*t2.Matrix[15],
		t.Matrix[3]*t2.Matrix[0]+t.Matrix[7]*t2.Matrix[1]+t.Matrix[15]*t2.Matrix[3],
		t.Matrix[3]*t2.Matrix[4]+t.Matrix[7]*t2.Matrix[5]+t.Matrix[15]*t2.Matrix[7],
		t.Matrix[3]*t2.Matrix[12]+t.Matrix[7]*t2.Matrix[13]+t.Matrix[15]*t2.Matrix[15])
}

func (t *Transform) TranslateXY(x, y float32) {
	translation := NewTransformFrom3x3(1, 0, x,
		0, 1, y,
		0, 0, 1)
	t.Combine(translation)
}

func (t *Transform) Translate(offset Vector2) {
	t.TranslateXY(offset.X, offset.Y)
}

func (t *Transform) Rotate(angle float32) {
	rad := angle * math.Pi / 180.0
	cos := float32(math.Cos(float64(rad)))
	sin := float32(math.Sin(float64(rad)))

	rotation := NewTransformFrom3x3(cos, -sin, 0,
		sin, cos, 0,
		0, 0, 1)

	t.Combine(rotation)
}

func (t *Transform) RotateAboutXY(angle, centerX, centerY float32) {
	rad := angle * math.Pi / 180.0
	cos := float32(math.Cos(float64(rad)))
	sin := float32(math.Sin(float64(rad)))

	rotation := NewTransformFrom3x3(cos, -sin, centerX*(1-cos)+centerY*sin,
		sin, cos, centerY*(1-cos)-centerX*sin,
		0, 0, 1)

	t.Combine(rotation)
}

func (t *Transform) RotateAbout(angle float32, center Vector2) {
	t.RotateAboutXY(angle, center.X, center.Y)
}

func (t *Transform) ScaleXY(scaleX, scaleY float32) {
	scaling := NewTransformFrom3x3(scaleX, 0, 0, 0, scaleY, 0, 0, 0, 1)
	t.Combine(scaling)
}

func (t *Transform) Scale(s Vector2) {
	t.ScaleXY(s.X, s.Y)
}

func (t *Transform) ScaleAboutXY(scaleX, scaleY, centerX, centerY float32) {
	scaling := NewTransformFrom3x3(scaleX, 0, centerX*(1-scaleX),
		0, scaleY, centerY*(1-scaleY),
		0, 0, 1)

	t.Combine(scaling)
}

func (t *Transform) ScaleAbout(factors, center Vector2) {
	t.ScaleAboutXY(factors.X, factors.Y, center.X, center.Y)
}
