// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program demonstrates the use of a TextureAtlas to tightly pack
// a number of small images into a single texture.
//
// It is an implementation of the 'Skyline-Bottom-Left' packing algorithm.
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glh"
	"github.com/go-gl/glu"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
)

// We'll use a square power-of-two texture.
// This is not strictly necessary, but offers better compatibility
// with older graphics drivers.
const AtlasSize = 512

func main() {
	err := initGL()
	if err != nil {
		log.Printf("InitGL: %v", err)
		return
	}

	defer glfw.Terminate()

	// Create our texture atlas.
	atlas := glh.NewTextureAtlas(AtlasSize, AtlasSize, 4)
	defer atlas.Release()

	// Fill the altas with image data.
	fillAtlas(atlas)

	// Display the atlas texture on a quad, so we can see
	// what it looks like.
	for glfw.WindowParam(glfw.Opened) > 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Bind the atlas texture and render a quad with it.
		atlas.Bind()
		gl.Begin(gl.QUADS)
		gl.TexCoord2f(0, 0)
		gl.Vertex2f(0, 0)
		gl.TexCoord2f(1, 0)
		gl.Vertex2f(AtlasSize, 0)
		gl.TexCoord2f(1, 1)
		gl.Vertex2f(AtlasSize, AtlasSize)
		gl.TexCoord2f(0, 1)
		gl.Vertex2f(0, AtlasSize)
		gl.End()
		atlas.Unbind()

		glfw.SwapBuffers()
	}
}

// fillAtlas fills the atlas with randomly sized and coloured rects.
func fillAtlas(a *glh.TextureAtlas) {
	// Create RNG with the same seed.
	rnd := rand.New(rand.NewSource(1e9))

	// generate a sequence of pseudo-randomly sized and coloured rectangles.
	// These represent the images we want to store in the atlas.
	for i := 0; i < 512; i++ {
		// Define the size of our image.
		size := int(rnd.Int31n(20)) + 10

		// Request a new region allocation from the atlas.
		// This region defines where our image is going to be stored
		// in the texture.
		region, ok := a.Allocate(size, size)

		// This will be false if the requested space could not be found
		// in the atlas. This means that either the requested size exceeds
		// the atlas bounds, or the atlas is full.
		if !ok {
			log.Printf("%d: No room for %dx%d", i, size, size)
			continue
		}

		// Create an image of the chosen size and fill it with a randomly
		// chosen, uniform color.
		rect := image.Rect(0, 0, region.W, region.H)
		img := image.NewRGBA(rect)

		clr := color.RGBA{
			byte(rnd.Int31n(200)) + 55,
			byte(rnd.Int31n(200)) + 55,
			byte(rnd.Int31n(200)) + 55,
			255,
		}

		draw.Draw(img, rect, image.NewUniform(clr), image.ZP, draw.Src)

		// Save the image to our atlas in the previously allocated region.
		a.Set(region, img.Pix, rect.Dx()*a.Depth())
	}

	// Save the texture as a PNG.
	//atlas.Save("atlas.png")

	// This creates the texture from the atlas pixel data.
	a.Commit()
}

// initGL initializes GLFW and OpenGL.
func initGL() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	err = glfw.OpenWindow(AtlasSize, AtlasSize, 8, 8, 8, 8, 0, 0, glfw.Windowed)
	if err != nil {
		glfw.Terminate()
		return err
	}

	glfw.SetWindowTitle("go-gl/gltext: freetype-gl")
	glfw.SetSwapInterval(1)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(onKey)

	errno := gl.Init()
	if errno != gl.NO_ERROR {
		str, err := glu.ErrorString(errno)
		if err != nil {
			return fmt.Errorf("Unknown openGL error: %d", errno)
		}
		return fmt.Errorf(str)
	}

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.2, 0.2, 0.23, 0.0)
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
