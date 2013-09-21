package sf

import (
	"testing"
)

func TestInverse(t *testing.T) {
	t1 := IdentityTransform()
	t1.Translate(Vector2{5, 0})
	t1.Rotate(180)
	t1.ScaleXY(2, 2)

	t2 := t1.Inverse()
	t1.Combine(t2)
	if t1 != IdentityTransform() {
		t.Fail()
	}
}
