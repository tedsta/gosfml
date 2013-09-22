package sf

type Rect struct {
	Left float32
	Top  float32
	W    float32
	H    float32
}

func (r Rect) IntersectsWith(r2 Rect) bool {
	if r.Left <= r2.Left+r2.W && r.Left+r.W >= r2.Left && r.Top <= r2.Top+r2.H && r.Top+r.H >= r2.Top {
		return true
	}
	return false
}
