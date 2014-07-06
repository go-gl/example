// Copyright 2013 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Draw a triangle using modern OpenGL. Press Esc to exit.
package main

import (
	"fmt"
	"github.com/andrebq/gas"
	gl "github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	glh "github.com/go-gl/glh"
	"io/ioutil"
)

const (
	Title    = "triangle"
	Width    = 500
	Height   = 500
	DataPath = "github.com/go-gl/examples/data"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

// Load file from the data directory and return its contents or panic on error.
func loadDataFile(filePath string) string {
	absFilePath, err := gas.Abs(DataPath + "/" + filePath)
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("glfw init failed")
	}
	defer glfw.Terminate()

	// use OpenGL 3.3 with deprecated functionality removed
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(keyCallback)
	window.MakeContextCurrent()

	// use vsync
	glfw.SwapInterval(1)

	// have to call this here or some OpenGL calls like CreateProgram will cause segfault
	if gl.Init() != 0 {
		panic("glew init failed")
	}
	gl.GetError() // ignore INVALID_ENUM that GLEW raises when using OpenGL 3.2+

	vShader := glh.Shader{gl.VERTEX_SHADER, loadDataFile("triangle.v.glsl")}
	fShader := glh.Shader{gl.FRAGMENT_SHADER, loadDataFile("triangle.f.glsl")}
	program := glh.NewProgram(vShader, fShader)
	program.Use()

	// OpenGL requires us to have a vertex array object bound
	vertexArray := gl.GenVertexArray()
	vertexArray.Bind()

	triangleVertices := [...]float32{-0.5, -0.5, -0.5, 0.5, 0.5, -0.5}
	triangleBuffer := gl.GenBuffer()
	triangleBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(glh.Sizeof(gl.FLOAT))*len(triangleVertices), &triangleVertices, gl.STATIC_DRAW)

	positionLocation := program.GetAttribLocation("position")
	positionLocation.EnableArray()
	positionLocation.AttribPointer(2, gl.FLOAT, false, 0, nil)

	for !window.ShouldClose() {
		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, width, height)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, len(triangleVertices))

		window.SwapBuffers()
		glfw.PollEvents()
	}

	// tell us if there were any errors
	glh.OpenGLSentinel()
}
