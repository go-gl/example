// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This test opens a window with FSAA sampling enabled, then verifies that we
// indeed got a window with > 0 sampling enabled. Creation of the window is
// necessary to make the test work.
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glu"
	"os"
)

const GL_MULTISAMPLE_ARB = 0x809D

func main() {
	var err error
	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	defer glfw.Terminate()

	// Open window with FSAA samples (if possible).
	glfw.OpenWindowHint(glfw.FsaaSamples, 4)

	if err = glfw.OpenWindow(400, 400, 0, 0, 0, 0, 0, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	defer glfw.CloseWindow()

	glfw.SetWindowTitle("Aliasing Detector")
	glfw.SetSwapInterval(1)

	if samples := glfw.WindowParam(glfw.FsaaSamples); samples != 0 {
		fmt.Fprintf(os.Stdout, "Context reports FSAA is supported with %d samples\n", samples)
	} else {
		fmt.Fprintf(os.Stdout, "Context reports FSAA is unsupported\n")
	}

	gl.MatrixMode(gl.PROJECTION)
	glu.Perspective(0, 1, 0, 1)

	for glfw.WindowParam(glfw.Opened) == 1 {
		time := float32(glfw.Time())

		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.LoadIdentity()
		gl.Translatef(0.5, 0, 0)
		gl.Rotatef(time, 0, 0, 1)

		gl.Enable(GL_MULTISAMPLE_ARB)
		gl.Color3f(1, 1, 1)
		gl.Rectf(-0.25, -0.25, 0.25, 0.25)

		gl.LoadIdentity()
		gl.Translatef(-0.5, 0, 0)
		gl.Rotatef(time, 0, 0, 1)

		gl.Disable(GL_MULTISAMPLE_ARB)
		gl.Color3f(1, 1, 1)
		gl.Rectf(-0.25, -0.25, 0.25, 0.25)
		glfw.SwapBuffers()
	}
}
