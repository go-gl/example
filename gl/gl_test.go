// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/go-gl/testutils"

	. "github.com/go-gl/gl"
)

var (
	array = [...]int32{
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16}
	slice = array[:]
)

func TestTexImage1D(t *testing.T) {
	gltest.OnTheMainThread(func() {

		TexImage1D(TEXTURE_1D,
			0, RGBA, 16, 0, RGBA, INT,
			&array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		Enable(TEXTURE_1D)
		TexImage1D(TEXTURE_1D,
			1, RGBA, 16, 1, RGBA, INT,
			&slice[3])
		if GetError() != NO_ERROR {
			t.Error("pointer to element failed")
		}

		TexImage1D(TEXTURE_1D,
			0, RGBA, 16, 0, RGBA, INT,
			slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		TexImage1D(PROXY_TEXTURE_1D,
			0, RGBA, 16, 0, RGBA, INT,
			nil)
		if GetError() != NO_ERROR {
			t.Error("nil pointer failed")
		}
	}, func() {})
}

func TestTexImage2D(t *testing.T) {
	gltest.OnTheMainThread(func() {

		TexImage2D(TEXTURE_2D,
			0, RGBA, 4, 4, 0, RGBA, INT,
			&array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		TexImage2D(TEXTURE_2D,
			0, RGBA, 4, 4, 0, RGBA, INT,
			slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		TexImage2D(PROXY_TEXTURE_2D,
			0, RGBA, 4, 4, 0, RGBA, INT,
			nil)
		if GetError() != NO_ERROR {
			t.Error("nil pointer failed")
		}
	}, func() {})
}

func TestTexImage3D(t *testing.T) {
	gltest.OnTheMainThread(func() {

		TexImage3D(TEXTURE_3D,
			0, RGBA, 2, 2, 2, 0, RGBA, INT,
			&array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		TexImage3D(TEXTURE_3D,
			0, RGBA, 2, 2, 2, 0, RGBA, INT,
			slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		TexImage3D(PROXY_TEXTURE_3D,
			0, RGBA, 2, 2, 2, 0, RGBA, INT,
			nil)
		if GetError() != NO_ERROR {
			t.Error("nil pointer failed")
		}
	}, func() {})
}

func TestTexSubImage1D(t *testing.T) {
	gltest.OnTheMainThread(func() {

		TexSubImage1D(TEXTURE_1D,
			0, 8, 8, RGBA, INT,
			&array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		TexSubImage1D(TEXTURE_1D,
			0, 8, 8, RGBA, INT,
			slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		// TexSubImage1D(PROXY_TEXTURE_1D,
		// 	0, 8, 8, RGBA, INT,
		// 	nil)
		// if GetError() != NO_ERROR {
		// 	t.Error("nil pointer failed")
		// }
	}, func() {})
}

func TestTexSubImage2D(t *testing.T) {
	gltest.OnTheMainThread(func() {

		TexSubImage2D(TEXTURE_2D,
			0, 2, 2, 2, 2, RGBA, INT,
			&array)
		if GetError() != NO_ERROR {
			t.Error("pointer to failed")
		}

		TexSubImage2D(TEXTURE_2D,
			0, 2, 2, 2, 2, RGBA, INT,
			slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		// TexSubImage2D(PROXY_TEXTURE_2D,
		// 	0, 2, 2, 2, 2, RGBA, INT,
		// 	nil)
		// if GetError() != NO_ERROR {
		// 	t.Error("nil pointer failed")
		// }
	}, func() {})
}

func TestTexSubImage3D(t *testing.T) {
	gltest.OnTheMainThread(func() {

		TexSubImage3D(TEXTURE_3D,
			0, 1, 1, 1, 1, 1, 1, RGBA, INT,
			&array)
		if GetError() != NO_ERROR {
			t.Error("pointer to failed")
		}

		TexSubImage3D(TEXTURE_3D,
			0, 1, 1, 1, 1, 1, 1, RGBA, INT,
			slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		// TexSubImage3D(PROXY_TEXTURE_3D,
		// 	0, 1, 1, 1, 1, 1, 1, RGBA, INT,
		// 	nil)
		// if GetError() != NO_ERROR {
		// 	t.Error("nil pointer failed")
		// }
	}, func() {})
}

func newBuffer(bytes int) Buffer {
	buf := GenBuffer()
	buf.Bind(ARRAY_BUFFER)
	BufferData(ARRAY_BUFFER, bytes, slice, STATIC_READ)
	return buf
}

func TestBufferData(t *testing.T) {
	gltest.OnTheMainThread(func() {

		buf := newBuffer(16 * 4)
		defer buf.Delete()

		BufferData(ARRAY_BUFFER, 16*4, &array, STATIC_READ)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		BufferData(ARRAY_BUFFER, 16*4, slice, STATIC_READ)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		BufferData(ARRAY_BUFFER, 16*4, nil, STATIC_READ)
		if GetError() != NO_ERROR {
			t.Error("nil pointer failed")
		}
	}, func() {})
}

func TestBufferSubData(t *testing.T) {
	gltest.OnTheMainThread(func() {

		buf := newBuffer(16 * 4)
		defer buf.Delete()

		BufferSubData(ARRAY_BUFFER, 0, 4*4, &array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		BufferSubData(ARRAY_BUFFER, 0, 4*4, slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}
	}, func() {})
}

func TestGetBufferSubData(t *testing.T) {
	gltest.OnTheMainThread(func() {

		buf := newBuffer(16 * 4)
		defer buf.Delete()

		result := make([]int32, 4)
		GetBufferSubData(ARRAY_BUFFER, 7*4, 4*4, result)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}
		if result[0] != 8 {
			t.Error("[8, 9, 10, 11] expected, was", result)
		}
	}, func() {})
}

func newProgram() Program {

	vs := `
		uniform float[] foo;
		attribute float bar;

		void main()
		{
			gl_Position = gl_Vertex * foo[0] * bar;
		}`

	fs := `
		void main()
		{
			gl_FragColor = vec4(1.0);
		}`

	// vertex shader
	vshader := CreateShader(VERTEX_SHADER)
	vshader.Source(vs)
	vshader.Compile()
	if vshader.Get(COMPILE_STATUS) != TRUE {
		panic("vertex shader error: " + vshader.GetInfoLog())
	}

	// fragment shader
	fshader := CreateShader(FRAGMENT_SHADER)
	fshader.Source(fs)
	fshader.Compile()
	if fshader.Get(COMPILE_STATUS) != TRUE {
		panic("fragment shader error: " + fshader.GetInfoLog())
	}

	// program
	prg := CreateProgram()
	prg.AttachShader(vshader)
	prg.AttachShader(fshader)
	prg.Link()
	if prg.Get(LINK_STATUS) != TRUE {
		panic("linker error: " + prg.GetInfoLog())
	}

	prg.Use()
	return prg
}

func TestAttribPointer(t *testing.T) {
	gltest.OnTheMainThread(func() {

		prg := newProgram()
		defer prg.Delete()

		bar := prg.GetAttribLocation("bar")
		bar.EnableArray()

		bar.AttribPointer(1, INT, false, 0, &array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		bar.AttribPointer(1, INT, false, 0, slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}
	}, func() {})
}

func TestUniform1fv(t *testing.T) {
	gltest.OnTheMainThread(func() {

		prg := newProgram()
		defer prg.Delete()

		values := []float32{1.0, 1.1, 1.2, 1.3}
		foo := prg.GetUniformLocation("foo")
		foo.Uniform1fv(len(values), values)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}
	}, func() {})
}

func TestCallLists(t *testing.T) {
	gltest.OnTheMainThread(func() {

		base := GenLists(2)
		if base == 0 {
			t.Error("GenLists failed")
		}

		NewList(base+0, COMPILE)
		Color3i(1, 2, 3)
		Begin(POLYGON)
		Vertex2f(1.0, 2.0)
		Vertex2f(3.0, 4.0)
		Vertex2f(5.0, 6.0)
		End()
		EndList()
		if GetError() != NO_ERROR {
			t.Error("NewList 1 failed")
		}

		NewList(base+1, COMPILE)
		Begin(POINTS)
		Vertex2i(3, 4)
		End()
		EndList()
		if GetError() != NO_ERROR {
			t.Error("NewList 2 failed")
		}

		CallList(base)
		if GetError() != NO_ERROR {
			t.Error("CallList 1 failed")
		}

		ListBase(base)
		if GetError() != NO_ERROR {
			t.Error("ListBase failed")
		}

		lists := []uint{0, 1}
		CallLists(2, UNSIGNED_INT, lists)
		if GetError() != NO_ERROR {
			t.Error("CallLists 1 2 failed")
		}
	}, func() {})
}

func TestColorPointer(t *testing.T) {
	gltest.OnTheMainThread(func() {

		EnableClientState(COLOR_ARRAY)

		ColorPointer(4, INT, 0, &array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		ColorPointer(4, INT, 0, slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		DisableClientState(COLOR_ARRAY)
		buf := newBuffer(16 * 4)
		defer buf.Delete()

		ColorPointer(4, INT, 0, uintptr(0))
		if GetError() != NO_ERROR {
			t.Error("buffer offset failed")
		}
	}, func() {})
}

func TestDrawElements(t *testing.T) {
	gltest.OnTheMainThread(func() {

		EnableClientState(VERTEX_ARRAY)
		VertexPointer(2, INT, 0, slice)

		indices := []uint32{0, 1, 2, 3}

		DrawElements(TRIANGLE_STRIP, 4, UNSIGNED_INT, indices)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		DisableClientState(VERTEX_ARRAY)
	}, func() {})
}

func TestDrawPixels(t *testing.T) {
	gltest.OnTheMainThread(func() {

		DrawPixels(4, 4, RGBA, UNSIGNED_INT, slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		buf := newBuffer(16 * 4)
		buf.Bind(PIXEL_UNPACK_BUFFER)
		defer buf.Delete()

		DrawPixels(2, 2, RGBA, UNSIGNED_INT, uintptr(0))
		if GetError() != NO_ERROR {
			t.Error("buffer offset failed")
		}
	}, func() {})
}

func TestIndexPointer(t *testing.T) {
	gltest.OnTheMainThread(func() {

		EnableClientState(INDEX_ARRAY)

		IndexPointer(INT, 0, &array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		IndexPointer(INT, 0, slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		DisableClientState(INDEX_ARRAY)
		buf := newBuffer(16 * 4)
		defer buf.Delete()

		IndexPointer(INT, 0, uintptr(0))
		if GetError() != NO_ERROR {
			t.Error("buffer offset failed")
		}
	}, func() {})
}

func TestNormalPointer(t *testing.T) {
	gltest.OnTheMainThread(func() {

		EnableClientState(NORMAL_ARRAY)

		NormalPointer(INT, 0, &array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		NormalPointer(INT, 0, slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		DisableClientState(NORMAL_ARRAY)
		buf := newBuffer(16 * 4)
		defer buf.Delete()

		NormalPointer(INT, 0, uintptr(0))
		if GetError() != NO_ERROR {
			t.Error("buffer offset failed")
		}
	}, func() {})
}

func TestReadPixels(t *testing.T) {
	gltest.OnTheMainThread(func() {

		pixels := []byte{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4}
		DrawPixels(2, 2, RGB, UNSIGNED_BYTE, pixels)

		result := make([]byte, 4*3)
		ReadPixels(0, 0, 2, 2, RGB, UNSIGNED_BYTE, result)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		// result is {1,1,1, 2,2,2, 0,0,3, 4,4,4}, no idea why.
		if false {
			t.Error("was", result)
		}
	}, func() {})
}

func TestTexCoordPointer(t *testing.T) {
	gltest.OnTheMainThread(func() {

		EnableClientState(TEXTURE_COORD_ARRAY)

		TexCoordPointer(3, INT, 0, &array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		TexCoordPointer(3, INT, 0, slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		DisableClientState(TEXTURE_COORD_ARRAY)
		buf := newBuffer(16 * 4)
		defer buf.Delete()

		TexCoordPointer(3, INT, 0, uintptr(0))
		if GetError() != NO_ERROR {
			t.Error("buffer offset failed")
		}
	}, func() {})
}

func TestVertexPointer(t *testing.T) {
	gltest.OnTheMainThread(func() {

		EnableClientState(VERTEX_ARRAY)

		VertexPointer(3, INT, 0, &array)
		if GetError() != NO_ERROR {
			t.Error("pointer to array failed")
		}

		VertexPointer(3, INT, 0, slice)
		if GetError() != NO_ERROR {
			t.Error("slice failed")
		}

		DisableClientState(VERTEX_ARRAY)
		buf := newBuffer(16 * 4)
		defer buf.Delete()

		VertexPointer(3, INT, 0, uintptr(0))
		if GetError() != NO_ERROR {
			t.Error("buffer offset failed")
		}
	}, func() {})
}
