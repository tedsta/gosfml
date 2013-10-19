package sf

import (
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

func init() {
	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	gl.Init()
}
