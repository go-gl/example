// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program demonstrates the use of truetype font rendering.
// It renders the same string using different font scales.
package main

import (
	"fmt"
	"github.com/andrebq/gas"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/gltext"
	"github.com/go-gl/glu"
	"log"
	"os"
)

const SampleString = "0 1 2 3 4 5 6 7 8 9 A B C D E F"

var fonts [16]*gltext.Font

func main() {
	file, err := gas.Abs("code.google.com/p/freetype-go/testdata/luxisr.ttf")
	if err != nil {
		log.Printf("Find font file: %v", err)
		return
	}

	err = initGL()
	if err != nil {
		log.Printf("InitGL: %v", err)
		return
	}

	defer glfw.Terminate()

	// Load the same font at different scale factors and directions.
	for i := range fonts {
		fonts[i], err = loadFont(file, int32(12+i))
		if err != nil {
			log.Printf("LoadFont: %v", err)
			return
		}

		defer fonts[i].Release()
	}

	for glfw.WindowParam(glfw.Opened) > 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		err = drawString(10, 10, SampleString)
		if err != nil {
			log.Printf("Printf: %v", err)
			return
		}

		glfw.SwapBuffers()
	}
}

// loadFont loads the specified font at the given scale.
func loadFont(file string, scale int32) (*gltext.Font, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	return gltext.LoadTruetype(fd, scale, 32, 127, gltext.LeftToRight)
}

// drawString draws the same string for each loaded font.
func drawString(x, y float32, str string) error {
	for i := range fonts {
		if fonts[i] == nil {
			continue
		}

		// We need to offset each string by the height of the
		// font. To ensure they don't overlap each other.
		_, h := fonts[i].GlyphBounds()
		y := y + float32(i*h)

		// Draw a rectangular backdrop using the string's metrics.
		sw, sh := fonts[i].Metrics(SampleString)
		gl.Color4f(0.1, 0.1, 0.1, 0.7)
		gl.Rectf(x, y, x+float32(sw), y+float32(sh))

		// Render the string.
		gl.Color4f(1, 1, 1, 1)
		err := fonts[i].Printf(x, y, str)
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

	glfw.SetWindowTitle("go-gl/gltext: Truetype font example")
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
