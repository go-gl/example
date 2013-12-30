// Copyright 2013 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Draw a triangle using modern OpenGL. Press Esc to exit.
package main

import (
	"fmt"
	gl "github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	glh "github.com/go-gl/glh"
	"os"
)

const (
	Title  = "triangle"
	Width  = 500
	Height = 500
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

// Return file's contents or panic on error.
func loadFile(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	b := make([]byte, info.Size())
	if _, err := file.Read(b); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}
	return b
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

	vShader := glh.Shader{gl.VERTEX_SHADER, string(loadFile("../../data/triangle.v.glsl"))}
	fShader := glh.Shader{gl.FRAGMENT_SHADER, string(loadFile("../../data/triangle.f.glsl"))}
	program := glh.NewProgram(vShader, fShader)
	program.Use()

	// OpenGL requires us to have a vertex array object bound
	vertexArray := gl.GenVertexArray()
	vertexArray.Bind()

	triangleVertices := [...]float32{-0.5, -0.5, -0.5, 0.5, 0.5, -0.5} // float32 should be equivalent to GLfloat
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
