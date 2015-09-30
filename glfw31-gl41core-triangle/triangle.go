package main

import (
	"fmt"
	"runtime"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	glfw "github.com/go-gl/glfw/v3.1/glfw"
)

const windowWidth = 800
const windowHeight = 640

var vertexShader = `
#version 410

in vec3 pos;
out vec3 position; // used in fragment shader

void main() {
  position = pos;
  gl_Position = vec4(pos, 1.0);
}
` + "\x00"

var fragmentShader = `
#version 410

in vec3 position; // comes from vertex shader
out vec4 outColor;

void main() {
  outColor = vec4(position+0.5, 1.0); // adds 0.5 to make it pretty (shifts the colors)
}
` + "\x00"

func main() {
	runtime.LockOSThread()

	initGlfw()
	setOpenGlVersion()
	var win = makeWindow()
	printOpenGlVersionInfo()

	var vao = createVao()

	gl.Viewport(0, 0, windowWidth, windowHeight)

	shader := createProgram()

	drawLoop(win, vao, shader)
}

func drawLoop(win *glfw.Window, vao uint32, shader uint32) {
	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.BindVertexArray(vao)
		gl.UseProgram(shader)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		glfw.PollEvents()
		win.SwapBuffers()
	}
}

func getPoints() []float32 {
	var points = []float32{
		-0.75, -0.75, 0,
		0.75, -0.75, 0,
		-0.75, 0.75, 0,
	}
	return points
}

func createVbo(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(&points[0]), gl.STATIC_DRAW)
	return vbo
}

func createVao() uint32 {
	var vbo = createVbo(getPoints())

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return vao
}

func setOpenGlVersion() {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)    // Necessary for OS X
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) // Necessary for OS X
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	glfw.WindowHint(glfw.Resizable, glfw.True)
}

func createProgram() uint32 {
	vs := gl.CreateShader(gl.VERTEX_SHADER)
	cvertexShader := gl.Str(vertexShader)
	gl.ShaderSource(vs, 1, &cvertexShader, nil)
	gl.CompileShader(vs)

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	cfragmentShader := gl.Str(fragmentShader)
	gl.ShaderSource(fs, 1, &cfragmentShader, nil)
	gl.CompileShader(fs)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, fs)
	gl.AttachShader(shaderProgram, vs)

	gl.LinkProgram(shaderProgram)

	return shaderProgram
}

func printOpenGlVersionInfo() {
	fmt.Printf("%s\n", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Printf("%s\n", gl.GoStr(gl.GetString(gl.VERSION)))
}

func initGlfw() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
}

func makeWindow() *glfw.Window {
	win, err := glfw.CreateWindow(windowWidth, windowHeight, "Tutorial #1", nil, nil)

	if err != nil {
		panic(err)
	}

	if err := gl.Init(); err != nil {
		panic(err)
	}

	win.MakeContextCurrent()
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	return win
}
