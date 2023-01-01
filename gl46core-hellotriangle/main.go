// Copyright 2021 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a red triangle using GLFW 3.3 and OpenGL 4.6 core forward-compatible profile.
package main

import (
	"errors"
	"fmt"
	_ "image/png"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Hello colorful triangle", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow.
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Compile, link and validate vertex and fragment shaders.
	program, err := CompileProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// float32 is 4 bytes wide.
	const attrSize = 4
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, attrSize*len(triangleVertices), gl.Ptr(triangleVertices), gl.STATIC_DRAW)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)

	// Stride is 2 since our data is 2D.
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 2*attrSize, 0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}

var vertexShader string = `
#version 330
in vec3 vert;

void main() {
	gl_Position = vec4(vert,1);
}
` + "\x00"

var fragmentShader string = `
#version 330
out vec4 outputColor;

void main() {
	outputColor = vec4(1.0,0.0,0.0,1.0);
}
` + "\x00"

var triangleVertices = []float32{
	-1.0, -1.0,
	1.0, -1.0,
	0.0, 1.0,
}

func CompileProgram(vertexSrcCode, fragmentSrcCode string) (program uint32, err error) {
	program = gl.CreateProgram()
	vid, err := compile(gl.VERTEX_SHADER, vertexSrcCode)
	if err != nil {
		return 0, fmt.Errorf("vertex shader compile: %w", err)
	}
	fid, err := compile(gl.FRAGMENT_SHADER, fragmentSrcCode)
	if err != nil {
		return 0, fmt.Errorf("fragment shader compile: %w", err)
	}
	gl.AttachShader(program, vid)
	gl.AttachShader(program, fid)
	gl.LinkProgram(program)
	log := ivLog(program, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog)
	if len(log) > 0 {
		return 0, fmt.Errorf("link failed: %v", log)
	}
	// We should technically call DetachShader after linking... https://www.youtube.com/watch?v=71BLZwRGUJE&list=PLlrATfBNZ98foTJPJ_Ev03o2oq3-GGOS2&index=7&ab_channel=TheCherno
	gl.ValidateProgram(program)
	log = ivLog(program, gl.VALIDATE_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog)
	if len(log) > 0 {
		return 0, fmt.Errorf("validation failed: %v", log)
	}

	// We can clean up.
	gl.DeleteShader(vid)
	gl.DeleteShader(fid)
	return program, nil
}

func compile(shaderType uint32, sourceCode string) (uint32, error) {
	id := gl.CreateShader(shaderType)
	csources, free := gl.Strs(sourceCode)
	gl.ShaderSource(id, 1, csources, nil)
	free()
	gl.CompileShader(id)

	// We now check the errors during compile, if there were any.
	log := ivLog(id, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog)
	if len(log) > 0 {
		return 0, errors.New(log)
	}
	return id, nil
}

// ivLog is a helper function for extracting log data
// from a Shader compilation step or program linking.
//
//	log := ivLog(id, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog)
//	if len(log) > 0 {
//		return 0, errors.New(log)
//	}
func ivLog(id, plName uint32, getIV func(program uint32, pname uint32, params *int32), getInfo func(program uint32, bufSize int32, length *int32, infoLog *uint8)) string {
	var iv int32
	getIV(id, plName, &iv)
	if iv == gl.FALSE {
		var logLength int32
		getIV(id, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength)
		getInfo(id, logLength, &logLength, &log[0])
		return string(log[:len(log)-1]) // we exclude the last null character.
	}
	return ""
}
