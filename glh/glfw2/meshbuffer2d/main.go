// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program demonstrates the use of a MeshBuffer.
// It displays a grid of quads with a randomly picked color.
// Hovering over each quad with the mouse, will assign it a new color.
//
// This shows how to modify mesh data after creating the buffer.
package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glh"
	"log"
	"math/rand"
)

const (
	PaletteSize = 6
	CellWidth   = 50
	CellHeight  = 50
	Cols        = 10
	Rows        = 10
)

var (
	selected = -1

	// Predefined color palette.
	colors = [PaletteSize][3]byte{
		{0xff, 0x99, 0},
		{0, 0xff, 0x99},
		{0x99, 0, 0xff},
		{0xff, 0, 0x99},
		{0x99, 0xff, 0},
		{0, 0x99, 0xff},
	}

	// RNG with fixed seed.
	rng = rand.New(rand.NewSource(1e9))
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

	attr := mb.Colors()

	// Perform the rendering.
	for glfw.WindowParam(glfw.Opened) > 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.LoadIdentity()

		// Center mesh on screen.
		wx, wy := glfw.WindowSize()
		px := (wx / 2) - ((Cols * CellWidth) / 2)
		py := (wy / 2) - ((Rows * CellHeight) / 2)
		gl.Translatef(float32(px), float32(py), 0)

		// Change the color of the quad under the mouse cursor.
		colorize(px, py, attr)

		// Render the mesh.
		mb.Render(gl.QUADS)

		glfw.SwapBuffers()
	}

	attr = nil
}

// createBuffer creates a mesh buffer and fills it with data.
func createBuffer() *glh.MeshBuffer {
	// Create a mesh buffer with the given attributes.
	mb := glh.NewMeshBuffer(
		glh.RenderBuffered,

		// 1 index per vertex.
		glh.NewIndexAttr(1, gl.UNSIGNED_SHORT, gl.STATIC_DRAW),

		// Vertex positions have 2 components (x, y).
		// These will never be changing, so mark them as static.
		glh.NewPositionAttr(2, gl.INT, gl.STATIC_DRAW),

		// Colors have 3 components (r, g, b).
		// These will be changing, so make them dynamic.
		glh.NewColorAttr(3, gl.UNSIGNED_BYTE, gl.DYNAMIC_DRAW),
	)

	// Define components for a simple, coloured quad.
	idx := []uint16{0, 1, 2, 3}
	pos := []int32{0, 0, 0, 0, 0, 0, 0, 0}
	clr := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

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

			palette := colors[rng.Int31n(6)]
			setColor(clr, palette[0], palette[1], palette[2])

			mb.Add(idx, pos, clr)
		}
	}

	return mb
}

// colorize changes the color of the quad, currently
// underneath the mouse cursor.
func colorize(px, py int, attr *glh.Attr) {
	// Fetch color data from mesh buffer.
	data := attr.Data().([]uint8)

	// We will be changing the color of the quad under the mouse cursor.
	// So first determine which quad the mouse is currently hovering over.
	mx, my := glfw.MousePos()
	mx = (mx - px) / CellWidth
	my = (my - py) / CellHeight
	index := my*Cols + mx

	// Ignore coordinates outside of the quad grid.
	// Make sure the index targets a valid quad.
	if mx < 0 || mx >= Cols || my < 0 || my >= Rows || index < 0 || index > Cols*Rows-1 {
		selected = -1
		return
	}

	if index == selected {
		return // No change.
	}

	// Set color for newly selected tile.
	if index != -1 {
		newcolor := colors[rng.Int31n(PaletteSize)]
		setColor(data[index*3*4:], newcolor[0], newcolor[1], newcolor[2])

		// Mark color data as stale, so it will be re-committed to the GPU
		// on the next Render call.
		//
		// This is technically only necessary for the RenderBuffer
		// and RenderShader modes, but it is good practice to always
		// use this call once data has changed.
		attr.Invalidate()
	}

	selected = index
}

// setColor sets the given color for a single quad.
//
// Each quad has four vertices. Each vertex has its own color with
// three components (r/g/b).
func setColor(d []byte, r, g, b byte) {
	d[0] = r
	d[1] = g
	d[2] = b

	d[3] = r
	d[4] = g
	d[5] = b

	d[6] = r
	d[7] = g
	d[8] = b

	d[9] = r
	d[10] = g
	d[11] = b
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
