// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

// Intentionally empty package here to make the package "go install"-able.

import (
	// These dependencies are listed here to prevent "go get" from having to
	// do more work later. Otherwise, the later delay would count against
	// the tests, which only have five minutes to run - network and build
	// activity can take a long time on a heavily loaded test node.
	_ "github.com/andrebq/gas"
	_ "github.com/banthar/Go-SDL/sdl"

	_ "github.com/go-gl/gl"
	_ "github.com/go-gl/glfw"
	_ "github.com/go-gl/gltext"
	_ "github.com/go-gl/glu"
	_ "github.com/go-gl/testutils"
)
