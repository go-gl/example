// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// NEHE Tutorial 06: Texture Mapping.
// http://nehe.gamedev.net/data/lessons/lesson.asp?lesson=06
package main

import (
	"errors"
	"log"

	"github.com/andrebq/gas"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glu"
)

const (
	Title  = "Nehe 06"
	Width  = 640
	Height = 480
)

var (
	running      bool
	rotation     [3]float32
	textures     []gl.Texture
	texturefiles [1]string
)

func init() {
	texturefiles[0], _ = gas.Abs("github.com/go-gl/examples/data/NeHe.tga")
}

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

	if err = initGL(); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	defer destroyGL()

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

func initGL() (err error) {
	if err = loadTextures(); err != nil {
		return
	}

	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_2D)
	return
}

func destroyGL() {
	gl.DeleteTextures(textures)
	textures = nil
}

func loadTextures() (err error) {
	textures = make([]gl.Texture, len(texturefiles))
	gl.GenTextures(textures)

	for i := range texturefiles {
		textures[i].Bind(gl.TEXTURE_2D)

		if !glfw.LoadTexture2D(texturefiles[i], 0) {
			return errors.New("Failed to load texture: " + texturefiles[i])
		}

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}

	return
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	gl.Translatef(0, 0, -5)

	gl.Rotatef(rotation[0], 1, 0, 0)
	gl.Rotatef(rotation[1], 0, 1, 0)
	gl.Rotatef(rotation[2], 0, 0, 1)

	rotation[0] += 0.3
	rotation[1] += 0.2
	rotation[2] += 0.4

	textures[0].Bind(gl.TEXTURE_2D)

	gl.Begin(gl.QUADS)
	// Front Face
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, 1) // Bottom Left Of The Texture and Quad
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, 1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, 1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, 1) // Top Left Of The Texture and Quad
	// Back Face
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, -1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, -1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, -1) // Top Left Of The Texture and Quad
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, -1) // Bottom Left Of The Texture and Quad
	// Top Face
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1) // Top Left Of The Texture and Quad
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, 1, 1) // Bottom Left Of The Texture and Quad
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, 1, 1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1) // Top Right Of The Texture and Quad
	// Bottom Face
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, -1, -1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, -1, -1) // Top Left Of The Texture and Quad
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1) // Bottom Left Of The Texture and Quad
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1) // Bottom Right Of The Texture and Quad
	// Right face
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, -1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, 1) // Top Left Of The Texture and Quad
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1) // Bottom Left Of The Texture and Quad
	// Left Face
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, -1) // Bottom Left Of The Texture and Quad
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, 1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1) // Top Left Of The Texture and Quad
	gl.End()

	glfw.SwapBuffers()
}
