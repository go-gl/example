// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// NEHE Tutorial 08: Blending.
// http://nehe.gamedev.net/data/lessons/lesson.asp?lesson=08
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
	Title  = "Nehe 08"
	Width  = 640
	Height = 480
)

var (
	running      bool
	textures     []gl.Texture = make([]gl.Texture, 3)
	texturefiles [1]string
	light        bool                                     // Display light?
	blend        bool                                     // Perform blending?
	rotation     [2]float32                               // X/Y rotation.
	speed        [2]float32                               // X/Y speed.
	z            float32    = -5                          // Depth into the scene.
	ambient      []float32  = []float32{0.5, 0.5, 0.5, 1} // ambient light colour.
	diffuse      []float32  = []float32{1, 1, 1, 1}       // diffuse light colour.
	lightpos     []float32  = []float32{0, 0, 2, 1}       // Position of light source.
	filter       int                                      // Index of current texture to display.
)

func init() {
	texturefiles[0], _ = gas.Abs("github.com/go-gl/examples/data/Glass.tga")
}

func main() {
	var err error
	if err = glfw.Init(); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	defer glfw.Terminate()

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 8, 0, glfw.Windowed); err != nil {
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
	case 'L':
		if state == 1 {
			if light = !light; !light {
				gl.Disable(gl.LIGHTING)
			} else {
				gl.Enable(gl.LIGHTING)
			}
		}
	case 'F':
		if state == 1 {
			if filter++; filter >= len(textures) {
				filter = 0
			}
		}
	case 'B': // B
		if state == 1 {
			if blend = !blend; blend {
				gl.Enable(gl.BLEND)
				gl.Disable(gl.DEPTH_TEST)
			} else {
				gl.Disable(gl.BLEND)
				gl.Enable(gl.DEPTH_TEST)
			}
		}
	case glfw.KeyPageup:
		z -= 0.2
	case glfw.KeyPagedown:
		z += 0.2
	case glfw.KeyUp:
		speed[0] -= 0.1
	case glfw.KeyDown:
		speed[0] += 0.1
	case glfw.KeyLeft:
		speed[1] -= 0.1
	case glfw.KeyRight:
		speed[1] += 0.1
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
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.DEPTH_TEST)
	gl.Color4f(1, 1, 1, 0.5)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	gl.Lightfv(gl.LIGHT1, gl.AMBIENT, ambient)
	gl.Lightfv(gl.LIGHT1, gl.AMBIENT, diffuse)
	gl.Lightfv(gl.LIGHT1, gl.POSITION, lightpos)
	gl.Enable(gl.LIGHT1)
	return
}

func destroyGL() {
	gl.DeleteTextures(textures)
	textures = nil
}

func loadTextures() (err error) {
	gl.GenTextures(textures)

	// Texture 1
	textures[0].Bind(gl.TEXTURE_2D)

	if !glfw.LoadTexture2D(texturefiles[0], 0) {
		return errors.New("Failed to load texture: " + texturefiles[0])
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// Texture 2
	textures[1].Bind(gl.TEXTURE_2D)

	if !glfw.LoadTexture2D(texturefiles[0], 0) {
		return errors.New("Failed to load texture: " + texturefiles[0])
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Texture 3
	textures[2].Bind(gl.TEXTURE_2D)

	if !glfw.LoadTexture2D(texturefiles[0], glfw.BuildMipmapsBit) {
		return errors.New("Failed to load texture: " + texturefiles[0])
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_NEAREST)

	return
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	gl.Translatef(0, 0, z)

	gl.Rotatef(rotation[0], 1, 0, 0)
	gl.Rotatef(rotation[1], 0, 1, 0)

	rotation[0] += speed[0]
	rotation[1] += speed[1]

	textures[filter].Bind(gl.TEXTURE_2D)

	gl.Begin(gl.QUADS)
	// Front Face
	gl.Normal3f(0, 0, 1) // Normal Pointing Towards Viewer
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, 1) // Point 1 (Front)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, 1) // Point 2 (Front)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, 1) // Point 3 (Front)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, 1) // Point 4 (Front)
	// Back Face
	gl.Normal3f(0, 0, -1) // Normal Pointing Away From Viewer
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, -1) // Point 1 (Back)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, -1) // Point 2 (Back)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, -1) // Point 3 (Back)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, -1) // Point 4 (Back)
	// Top Face
	gl.Normal3f(0, 1, 0) // Normal Pointing Up
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1) // Point 1 (Top)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, 1, 1) // Point 2 (Top)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, 1, 1) // Point 3 (Top)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1) // Point 4 (Top)
	// Bottom Face
	gl.Normal3f(0, -1, 0) // Normal Pointing Down
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, -1, -1) // Point 1 (Bottom)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, -1, -1) // Point 2 (Bottom)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1) // Point 3 (Bottom)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1) // Point 4 (Bottom)
	// Right face
	gl.Normal3f(1, 0, 0) // Normal Pointing Right
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, -1) // Point 1 (Right)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1) // Point 2 (Right)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, 1) // Point 3 (Right)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1) // Point 4 (Right)
	// Left Face
	gl.Normal3f(-1, 0, 0) // Normal Pointing Left
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, -1) // Point 1 (Left)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1) // Point 2 (Left)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, 1) // Point 3 (Left)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.End()

	glfw.SwapBuffers()
}
