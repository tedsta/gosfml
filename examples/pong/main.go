package main

import (
	"fmt"
	glfw "github.com/go-gl/glfw3"
	"github.com/tedsta/gosfml"
)

var (
	window           *glfw.Window
	target           *sf.RenderTarget
	p1, p2, ball     *Object
	p1Score, p2Score int
)

func main() {
	window, err := glfw.CreateWindow(800, 600, "Golang SFML Pong", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(onResize)
	window.SetKeyCallback(onKey)

	target = sf.NewRenderTarget(sf.Vector2{float32(800), float32(600)})

	p1 = NewObject(5, 5, 16, 64)
	p2 = NewObject(795-16, 5, 16, 64)
	ball = NewObject(400-8, 300-8, 16, 16)

	ball.vel = sf.Vector2{300, 300}

	clock := sf.NewClock()
	for !window.ShouldClose() {
		glfw.PollEvents()

		dt := float32(clock.Restart().Seconds())

		if ball.pos.X < 0 {
			ball.vel.X = 300
			p2Score++
			fmt.Println("Player 1: ", p1Score)
			fmt.Println("Player 2: ", p2Score)
			fmt.Println()
		}
		if ball.pos.X+ball.dim.X >= 800 {
			ball.vel.X = -300
			p1Score++
			fmt.Println("Player 1: ", p1Score)
			fmt.Println("Player 2: ", p2Score)
			fmt.Println()
		}
		if ball.pos.Y < 0 {
			ball.vel.Y = 300
		}
		if ball.pos.Y+ball.dim.Y >= 600 {
			ball.vel.Y = -300
		}

		if p1.Collision(ball) {
			ball.vel.X = 300
		}
		if p2.Collision(ball) {
			ball.vel.X = -300
		}

		p1.pos.X += p1.vel.X * dt
		p1.pos.Y += p1.vel.Y * dt
		p2.pos.X += p2.vel.X * dt
		p2.pos.Y += p2.vel.Y * dt
		ball.pos.X += ball.vel.X * dt
		ball.pos.Y += ball.vel.Y * dt

		// Render
		target.Clear(sf.Color{0, 0, 0, 0})
		p1.Render()
		p2.Render()
		ball.Render()
		window.SwapBuffers()
	}
}

// #############################################################################
// Object

type Object struct {
	pos sf.Vector2
	dim sf.Vector2
	vel sf.Vector2
}

func NewObject(x, y, w, h float32) *Object {
	return &Object{sf.Vector2{x, y}, sf.Vector2{w, h}, sf.Vector2{}}
}

func (o *Object) Collision(o2 *Object) bool {
	if o.pos.X+o.dim.X >= o2.pos.X && o.pos.X <= o2.pos.X+o2.dim.X &&
		o.pos.Y+o.dim.Y >= o2.pos.Y && o.pos.Y <= o2.pos.Y+o2.dim.Y {
		return true
	}
	return false
}

func (o *Object) Render() {
	var verts [4]sf.Vertex
	verts[0] = sf.Vertex{sf.Vector2{},
		sf.Color{255, 255, 255, 255},
		sf.Vector2{}}
	verts[1] = sf.Vertex{sf.Vector2{0, o.dim.Y},
		sf.Color{255, 255, 255, 255},
		sf.Vector2{0, o.dim.Y}}
	verts[2] = sf.Vertex{o.dim,
		sf.Color{255, 255, 255, 255},
		o.dim}
	verts[3] = sf.Vertex{sf.Vector2{o.dim.X, 0},
		sf.Color{255, 255, 255, 255},
		sf.Vector2{o.dim.X, 0}}

	states := sf.RenderStates{sf.BlendAlpha, sf.IdentityTransform(), nil}
	states.Transform.Translate(o.pos)
	target.Render(verts[:], sf.Quads, states)
}

// #############################################################################
// Callbacks

func onResize(wnd *glfw.Window, w, h int) {
	target.Size.X = float32(w)
	target.Size.Y = float32(h)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyW:
			p1.vel.Y = -300
		case glfw.KeyS:
			p1.vel.Y = 300
		case glfw.KeyUp:
			p2.vel.Y = -300
		case glfw.KeyDown:
			p2.vel.Y = 300
		}
	} else if action == glfw.Release {
		switch key {
		case glfw.KeyW, glfw.KeyS:
			p1.vel.Y = 0
		case glfw.KeyUp, glfw.KeyDown:
			p2.vel.Y = 0
		}
	}
}
