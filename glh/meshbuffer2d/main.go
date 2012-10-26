// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program demonstrates the use of a MeshBuffer.
package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glh"
	"log"
)

const (
	CellWidth  = 50
	CellHeight = 50
	Cols       = 10
	Rows       = 10
)

func main() {
	err := initGL()
	if err != nil {
		log.Printf("InitGL: %v", err)
		return
	}

	defer glfw.Terminate()

	mb := createBuffer()
	defer mb.Release()

	// Perform the rendering.
	for glfw.WindowParam(glfw.Opened) > 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.LoadIdentity()

		// Center mesh on screen.
		wx, wy := glfw.WindowSize()
		px := (wx / 2) - ((Cols * CellWidth) / 2)
		py := (wy / 2) - ((Rows * CellHeight) / 2)

		gl.Translatef(float32(px), float32(py), 0)
		mb.Render(gl.QUADS)

		glfw.SwapBuffers()
	}
}

func createBuffer() *glh.MeshBuffer {
	// Create a mesh buffer with the given attributes.
	mb := glh.NewMeshBuffer(
		// Indices have 1 element.
		glh.NewUint32Attr(1, gl.STATIC_DRAW),

		// Vertex positions have 2 components (x, y).
		glh.NewInt32Attr(2, gl.STATIC_DRAW),

		// Colors have 3 components (r, g, b).
		glh.NewFloat32Attr(3, gl.STATIC_DRAW),

		nil, // We have no surface normals.    
		nil, // No texture coordinates either. 
	)

	// Define components for a simple, coloured quad.
	idx := []uint32{0, 1, 2, 3}
	pos := []int32{0, 0, 0, 0, 0, 0, 0, 0}
	clr := []float32{1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1}

	for y := 0; y < Rows; y++ {
		py := int32(CellHeight * y)

		pos[1] = py + 1
		pos[3] = py + 1
		pos[5] = py + CellHeight
		pos[7] = py + CellHeight

		for x := 0; x < Cols; x++ {
			px := int32(CellWidth * x)
			pos[0] = px + 1
			pos[2] = px + CellWidth
			pos[4] = px + CellWidth
			pos[6] = px + 1

			mb.Add(idx, pos, clr, nil, nil)
		}
	}

	return mb
}

// initGL initializes GLFW and OpenGL.
func initGL() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.OpenWindowHint(glfw.FsaaSamples, 4)

	err = glfw.OpenWindow(512, 512, 8, 8, 8, 8, 0, 0, glfw.Windowed)
	if err != nil {
		glfw.Terminate()
		return err
	}

	glfw.SetWindowTitle("Meshbuffer 2D example")
	glfw.SetSwapInterval(1)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(onKey)

	gl.Init()
	if err = glh.CheckGLError(); err != nil {
		return err
	}

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)
	gl.Disable(gl.LIGHTING)
	gl.Enable(gl.COLOR_MATERIAL)
	gl.ClearColor(0.2, 0.2, 0.23, 1.0)
	return nil
}

// onKey handles key events.
func onKey(key, state int) {
	if key == glfw.KeyEsc {
		glfw.CloseWindow()
	}
}

// onResize handles window resize events.
func onResize(w, h int) {
	if w < 1 {
		w = 1
	}

	if h < 1 {
		h = 1
	}

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), float64(h), 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
