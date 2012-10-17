// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// NEHE Tutorial 04: Rotation.
// http://nehe.gamedev.net/data/lessons/lesson.asp?lesson=04
package main

import (
	"log"

	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glu"
)

const (
	Title  = "Nehe 04"
	Width  = 640
	Height = 480
)

var (
	trisAngle float32
	quadAngle float32
	running   bool
)

func main() {
	var err error
	if err = glfw.Init(); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	defer glfw.Terminate()

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 0, 8, glfw.Windowed); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(Title)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(onKey)

	initGL()

	running = true
	for running && glfw.WindowParam(glfw.Opened) == 1 {
		drawScene()
	}
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	glu.Perspective(45.0, float64(w)/float64(h), 0.1, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		running = false
	}
}

func initGL() {
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.LoadIdentity()
	gl.Translatef(-1.5, 0, -6)
	gl.Rotatef(trisAngle, 0, 1, 0)

	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1, 0, 0)
	gl.Vertex3f(0, 1, 0)
	gl.Color3f(0, 1, 0)
	gl.Vertex3f(-1, -1, 0)
	gl.Color3f(0, 0, 1)
	gl.Vertex3f(1, -1, 0)
	gl.End()

	gl.LoadIdentity()
	gl.Translatef(1.5, 0, -6)
	gl.Rotatef(quadAngle, 1, 0, 0)
	gl.Color3f(0.5, 0.5, 1.0)

	gl.Begin(gl.QUADS)
	gl.Vertex3f(-1, 1, 0)
	gl.Vertex3f(1, 1, 0)
	gl.Vertex3f(1, -1, 0)
	gl.Vertex3f(-1, -1, 0)
	gl.End()

	trisAngle += 0.2
	quadAngle -= 0.15

	glfw.SwapBuffers()
}
