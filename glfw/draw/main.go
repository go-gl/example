// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"math"

	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
)

var (
	running bool
	pen     Pen
	mouse   [3]int
)

func main() {
	var err error
	if err = glfw.Init(); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	defer glfw.Terminate()

	if err = glfw.OpenWindow(640, 480, 8, 8, 8, 8, 0, 0, glfw.Windowed); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	defer glfw.CloseWindow()

	glfw.SetWindowTitle("Draw")
	glfw.SetSwapInterval(1)
	glfw.SetKeyCallback(onKey)
	glfw.SetMouseButtonCallback(onMouseBtn)
	glfw.SetWindowSizeCallback(onResize)

	running = true
	for running && glfw.WindowParam(glfw.Opened) == 1 {
		if mouse[0] != 0 {
			pen.lineTo(glfw.MousePos())
		} else {
			pen.moveTo(glfw.MousePos())
		}

		glfw.SwapBuffers()
	}
}

func onResize(w, h int) {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, w, h)
	gl.Ortho(0, float64(w), float64(h), 0, -1.0, 1.0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func onMouseBtn(button, state int) {
	mouse[button] = state
}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		running = state == 0
	case 67: // 'c'
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}
}

func distanceTo(x1, y1, x2, y2 int) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return math.Sqrt(float64(dx*dx + dy*dy))
}

type Pen struct {
	pos    [2]int
	points [4096][2]int
	n      int
}

func (pen *Pen) lineTo(x, y int) {
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Color4f(0.0, 0.0, 0.0, 0.1)
	gl.Begin(gl.LINES)

	var p [2]int
	for i := range pen.points {
		p = pen.points[i]
		if p[0] == 0 && p[1] == 0 {
			continue
		}

		if distanceTo(x, y, p[0], p[1]) < 10.0 {
			gl.Vertex2i(x, y)
			gl.Vertex2i(p[0], p[1])
		}
	}

	gl.End()

	pen.n = (pen.n + 1) % len(pen.points)
	pen.points[pen.n][0] = x
	pen.points[pen.n][1] = y
	pen.moveTo(x, y)
}

func (pen *Pen) moveTo(x, y int) {
	pen.pos[0] = x
	pen.pos[1] = y
}
