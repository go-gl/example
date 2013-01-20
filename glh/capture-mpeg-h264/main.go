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
	"math"

	"image"
	"image/png"
	"os"

	"github.com/pwaller/go-ffmpeg-encoding"
)

func mkpng(filename string, im image.Image) {

	fd, err := os.Create(filename)
	if err != nil {
		log.Panic("Err: ", err)
	}
	defer fd.Close()

	png.Encode(fd, im)
}

func main() {
	err := initGL()
	if err != nil {
		log.Printf("InitGL: %v", err)
		return
	}

	defer glfw.Terminate()

	mb := createBuffer()
	defer mb.Release()

	gl.Enable(gl.BLEND)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	f, err := os.Create("test.mpeg")
	if err != nil {
		log.Panicf("Unable to open output file: %q", err)
	}

	im := image.NewRGBA(image.Rect(0, 0, 512, 512))
	e, err := ffmpeg.NewEncoder(ffmpeg.CODEC_ID_H264, im, f)
	if err != nil {
		log.Panicf("Unable to start encoder: %q", err)
	}

	// Perform the rendering.
	var angle float64
	for glfw.WindowParam(glfw.Opened) > 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.LoadIdentity()
		gl.Translatef(0, 0, -6)
		gl.Rotated(angle/10, 0, 1, 0)
		gl.Rotated(angle, 1, 1, 1)

		// Render a solid cube at half the scale.
		gl.Scalef(0.1, 0.1, 0.1)
		gl.Enable(gl.COLOR_MATERIAL)
		gl.Enable(gl.POLYGON_OFFSET_FILL)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		mb.Render(gl.QUADS)

		// Render wireframe cubes, with incremental size.
		gl.Disable(gl.COLOR_MATERIAL)
		gl.Disable(gl.POLYGON_OFFSET_FILL)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

		for i := 0; i < 1000; i++ {
			ifl := float64(i)
			scale := 0.000005*ifl + 1.0 + 0.01*math.Sin(angle/10+ifl/10)
			gl.Scaled(scale, scale, scale)
			gl.Rotated((ifl+angle)/1000, 1, 1, 1)
			gl.Rotated((ifl+angle)/1100, -1, 1, 1)

			mb.Render(gl.QUADS)
		}

		glh.ClearAlpha(1)
		glh.CaptureRGBA(im)

		err := e.WriteFrame()
		if err != nil {
			panic(err)
		}

		angle += 0.5
		//glh.CaptureToPng("test.png")
		glfw.SwapBuffers()
	}
	e.Close()
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

	a := float32(0.02)

	// Each vertex comes with its own colour.
	clr := []float32{
		1, 0, 0, a, 0, 1, 0, a, 0, 0, 1, a, 1, 0, 1, a,
		1, 1, 0, a, 0, 1, 1, a, 1, 1, 1, a, 0, 0, 0, a,
	}

	// These are the indices into the position and color lists.
	// They tell the GPU which position/color pair to use in order to construct
	// the whole cube. As can be seen, all elements are repeated multiple
	// times to create the correct layout. For large meshes with many duplicate
	// vertices, this can save a sizable amount of storage space.
	idx := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 3, 2, 5, 4,
		7, 6, 1, 0, 2, 1, 6, 5, 0, 3, 4, 7,
	}

	// Create a mesh buffer with the given attributes.
	mb := glh.NewMeshBuffer(
		glh.RenderBuffered,

		// Indices.
		glh.NewIndexAttr(1, gl.UNSIGNED_BYTE, gl.STATIC_DRAW),

		// Vertex positions have 3 components (x, y, z).
		glh.NewPositionAttr(3, gl.FLOAT, gl.STATIC_DRAW),

		// Colors have 4 components (r, g, b, a).
		glh.NewColorAttr(4, gl.FLOAT, gl.STATIC_DRAW),
	)

	// Add the mesh to the buffer.
	mb.Add(idx, pos, clr)
	return mb
}

// initGL initializes GLFW and OpenGL.
func initGL() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.OpenWindowHint(glfw.FsaaSamples, 4)
	glfw.OpenWindowHint(glfw.WindowNoResize, gl.TRUE)

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

	//gl.ClearColor(0.2, 0.2, 0.23, 1.0)
	gl.ClearColor(0, 0, 0, 1.0)
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
	glu.Perspective(45.0, float64(w)/float64(h), 0.1, 2000.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
