// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program demonstrates the use of bitmap (raster) fonts.
// It renders the same string using different font scale factors.
package main

import (
	"fmt"
	"github.com/andrebq/gas"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glu"
	"github.com/go-gl/text"
	"log"
	"os"
)

var fonts [8]*text.Font

func main() {
	// This file holds the actual glyph shapes.
	imgFile, err := gas.Abs("github.com/go-gl/examples/data/bitmap_font.png")
	if err != nil {
		log.Printf("Find font image file: %v", err)
		return
	}

	// This file is a JSON encoded dataset which describes the font
	// and contains the pixel offsets and sizes for each glyph in
	// bitmap_font.png. Both files are needed to load a bitmap font.
	configFile, err := gas.Abs("github.com/go-gl/examples/data/bitmap_font.js")
	if err != nil {
		log.Printf("Find font config file: %v", err)
		return
	}

	err = initGL()
	if err != nil {
		log.Printf("InitGL: %v", err)
		return
	}

	defer glfw.Terminate()

	// Load the same bitmap font at different scale factors.
	for i := range fonts {
		fonts[i], err = loadFont(imgFile, configFile, i+1)
		if err != nil {
			log.Printf("LoadFont: %v", err)
			return
		}

		defer fonts[i].Release()
	}

	for glfw.WindowParam(glfw.Opened) > 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		err = drawString(10, 10, "0 1 2 3 4 5 6 7 8 9 A B C D E F")
		if err != nil {
			log.Printf("Printf: %v", err)
			return
		}

		glfw.SwapBuffers()
	}
}

// loadFont loads a scaled bitmap font from the given image
// and configuration files.
func loadFont(image, config string, scale int) (*text.Font, error) {
	a, err := os.Open(image)
	if err != nil {
		return nil, err
	}

	defer a.Close()

	b, err := os.Open(config)
	if err != nil {
		return nil, err
	}

	defer b.Close()

	return text.LoadBitmap(a, b, scale)
}

// drawString draws the same string twice with a colour and location offset,
// to simulate a drop-shadow. It does so for each loaded font.
func drawString(x, y float32, str string) error {
	for i := range fonts {
		if fonts[i] == nil {
			continue
		}

		// We need to offset each string by the height of the
		// font. To ensure they don't overlap each other.
		_, h := fonts[i].GlyphBounds()
		y := y + float32(i*h)

		gl.Color4f(0.1, 0.1, 0.1, 0.7)
		err := fonts[i].Printf(x+2, y+2, str)
		if err != nil {
			return err
		}

		gl.Color4f(1, 1, 1, 1)
		err = fonts[i].Printf(x, y, str)
		if err != nil {
			return err
		}
	}

	return nil
}

// initGL initializes GLFW and OpenGL.
func initGL() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	err = glfw.OpenWindow(640, 480, 8, 8, 8, 8, 0, 0, glfw.Windowed)
	if err != nil {
		glfw.Terminate()
		return err
	}

	glfw.SetWindowTitle("go-gl/text: Bitmap font example")
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
