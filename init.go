package sf

import (
	"github.com/go-gl-legacy/gl"
	"github.com/go-gl/glfw3/v3.1/glfw"
)

func init() {
	if err := glfw.Init(); err != nil {
		panic("Can't init glfw!")
	}
	gl.Init()
}
