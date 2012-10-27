// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program demonstrates the use of a MeshBuffer.
package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glh"
	"github.com/go-gl/glu"
	"log"
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
	var angle float32
	for glfw.WindowParam(glfw.Opened) > 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.LoadIdentity()
		gl.Translatef(0, 0, -6)
		gl.Rotatef(angle, 1, 1, 1)

		// Render a solid cube at half the scale.
		gl.Scalef(0.2, 0.2, 0.2)
		gl.Enable(gl.COLOR_MATERIAL)
		gl.Enable(gl.POLYGON_OFFSET_FILL)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		mb.Render(gl.QUADS)

		// Render wireframe cubes, with incremental size.
		gl.Disable(gl.COLOR_MATERIAL)
		gl.Disable(gl.POLYGON_OFFSET_FILL)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

		for i := 0; i < 50; i++ {
			scale := 0.004*float32(i) + 1.0
			gl.Scalef(scale, scale, scale)
			mb.Render(gl.QUADS)
		}

		angle += 0.5
		glfw.SwapBuffers()
	}
}

func createBuffer() *glh.MeshBuffer {
	// We create as few vertices as possible.
	// Manually building a cube would require 24 vertices. Many of which
	// are duplicates. All we have to define here, is the 8 unique ones
	// necessary to construct each face of the cube.
	pos := []float32{
		1, 1, -1, -1, 1, -1, -1, 1, 1, 1, 1, 1,
		1, -1, 1, -1, -1, 1, -1, -1, -1, 1, -1, -1,
	}

	// Each vertex comes with its own colour.
	clr := []float32{
		1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 1,
		1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1,
	}

	// These are the indices into the Position and Color lists.
	// They tell the GPU which position/color pair to use in order to construct
	// the whole cube. As can be seen, all elements are repeated multiple
	// times to create the correct layout. For large meshes, this can save
	// a tremendous amount of storage space.
	idx := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 3, 2, 5, 4,
		7, 6, 1, 0, 2, 1, 6, 5, 0, 3, 4, 7,
	}

	// Create a mesh buffer with the given attributes.
	mb := glh.NewMeshBuffer(
		glh.RenderBuffered,

		// Indices.
		glh.NewMeshAttr(1, gl.UNSIGNED_BYTE, gl.STATIC_DRAW),

		// Vertex positions have 3 components (x, y, z).
		glh.NewMeshAttr(3, gl.FLOAT, gl.STATIC_DRAW),

		// Colors have 4 components (r, g, b, a).
		glh.NewMeshAttr(4, gl.FLOAT, gl.STATIC_DRAW),

		nil, // No surface normals.
		nil, // No texture coordinates.
	)

	// Add the mesh to the buffer.
	mb.Add(idx, pos, clr, nil, nil)
	return mb
}

// initGL initializes GLFW and OpenGL.
func initGL() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.OpenWindowHint(glfw.FsaaSamples, 4)

	err = glfw.OpenWindow(512, 512, 8, 8, 8, 8, 32, 0, glfw.Windowed)
	if err != nil {
		glfw.Terminate()
		return err
	}

	glfw.SetWindowTitle("Meshbuffer 3D example")
	glfw.SetSwapInterval(1)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(onKey)

	gl.Init()
	if err = glh.CheckGLError(); err != nil {
		return err
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)
	gl.Disable(gl.LIGHTING)

	gl.ClearColor(0.2, 0.2, 0.23, 1.0)
	gl.ShadeModel(gl.SMOOTH)
	gl.LineWidth(2)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
	gl.ColorMaterial(gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE)
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
	glu.Perspective(45.0, float64(w)/float64(h), 0.1, 200.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
