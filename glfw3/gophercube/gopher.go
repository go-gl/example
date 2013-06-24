//This example is copied from https://github.com/chsc/gogl/

package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

const (
	Title  = "Spinning Gopher"
	Width  = 640
	Height = 480
)

var (
	texture    gl.Texture
	rotx, roty float32
	ambient    []float32 = []float32{0.5, 0.5, 0.5, 1}
	diffuse    []float32 = []float32{1, 1, 1, 1}
	lightpos   []float32 = []float32{-5, 5, 10, 0}
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		panic(err)
	}
	
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	gl.Init()

	if err := initScene(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}
	defer destroyScene()

	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func createTexture(r io.Reader) (gl.Texture, error) {
	img, err := png.Decode(r)
	if err != nil {
		return gl.Texture(0), err
	}

	rgbaImg, ok := img.(*image.NRGBA)
	if !ok {
		return gl.Texture(0), errors.New("texture must be an NRGBA image")
	}

	textureId := gl.GenTexture()
	textureId.Bind(gl.TEXTURE_2D)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	// flip image: first pixel is lower left corner
	imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()
	data := make([]byte, imgWidth*imgHeight*4)
	lineLen := imgWidth * 4
	dest := len(data) - lineLen
	for src := 0; src < len(rgbaImg.Pix); src += rgbaImg.Stride {
		copy(data[dest:dest+lineLen], rgbaImg.Pix[src:src+rgbaImg.Stride])
		dest -= lineLen
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, 4, imgWidth, imgHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)

	return textureId, nil
}

func initScene() (err error) {
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)

	gl.ClearColor(0.5, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, ambient)
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, diffuse)
	gl.Lightfv(gl.LIGHT0, gl.POSITION, lightpos)
	gl.Enable(gl.LIGHT0)

	gl.Viewport(0, 0, Width, Height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	goph, err := os.Open("../../data/gopher.png")
	if err != nil {
		panic(err)
	}
	defer goph.Close()

	texture, err = createTexture(goph)
	return
}

func destroyScene() {
	texture.Delete()
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0, 0, -3.0)
	gl.Rotatef(rotx, 1, 0, 0)
	gl.Rotatef(roty, 0, 1, 0)

	rotx += 0.5
	roty += 0.5

	texture.Bind(gl.TEXTURE_2D)

	gl.Color4f(1, 1, 1, 1)

	gl.Begin(gl.QUADS)

	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, 1)

	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, -1)

	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)

	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)

	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)

	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)

	gl.End()
}
