// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	// These dependencies are listed here to prevent "go get" from having to
	// do more work later. Otherwise, the later delay would count against
	// the tests, which only have five minutes to run - network and build
	// activity can take a long time on a heavily loaded node.
	_ "github.com/andrebq/gas"
	_ "github.com/go-gl/gl"
	_ "github.com/go-gl/glfw"
	_ "github.com/go-gl/glu"
)

const default_sleeptime = 1 * time.Second

var runtimes = map[string]time.Duration{
	"nehe03": 5 * time.Second,
}

func runTest(t *testing.T, path string) {
	println(strings.Repeat("=", 80))
	println("-- subtest: ", path)
	test := exec.Command("go", "test", "-v", "./"+path)
	test.Stdout, test.Stderr = os.Stdout, os.Stderr

	err := test.Run()
	if err != nil {
		t.Fatal("Failed to run 'go test': ", err)
	}
}

func runExample(t *testing.T, path string, files []string) {
	println(strings.Repeat("=", 80))

	get := exec.Command("go", "get", "-d", "-v", "./"+path)
	get.Stdout = os.Stdout
	get.Stderr = os.Stderr
	err := get.Run()
	if err != nil {
		t.Fatal("Failed to run go get: ", err)
	}

	bin_name := filepath.Join("bin", path)
	bld := exec.Command("go", "build", "-v", "-o", bin_name, "./"+path)
	bld.Stdout = os.Stdout
	bld.Stderr = os.Stderr
	err = bld.Run()
	if err != nil {
		t.Fatal("Failed to run go build: ", err)
	}

	println(strings.Repeat("-", 80))

	cmd := exec.Command(bin_name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	done := false

	go func() {
		sleeptime, ok := runtimes[path]
		if !ok {
			sleeptime = default_sleeptime
		}
		time.Sleep(sleeptime)
		done = true

		defer func() {
			e := recover()
			if e == nil {
				return
			}
			err, ok := e.(error)
			if !ok {
				panic(e)
			}
			if err.Error() != "os: process already finished" {
				t.Fatal("Failed to terminate process: ", err)
			}
		}()
		err := cmd.Process.Kill()
		if err != nil {
			t.Fatal("Failed to terminate process: ", err)
		}
	}()

	err = cmd.Run()

	// If the done flag is true, then we made it through five seconds of runtime
	if !done && err != nil {
		//panic(err)
		t.Fatal("Process died unexpectedly: ", err)
	}
}

func goInstall() {
	install := exec.Command("go", "install", "-v")
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	err := install.Run()
	if err != nil {
		panic(err)
	}
}

func hasTest(files []string) bool {
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			return true
		}
	}
	return false
}

func hasNonTest(files []string) bool {
	for _, f := range files {
		if !strings.HasSuffix(f, "_test.go") {
			return true
		}
	}
	return false
}

// Runs all examples in turn using "go build" and "dirname/dirname".
// They run in an arbitrary order.
func TestExamples(t *testing.T) {
	// Required to ensure that the test assets are available at the right path.
	// Maybe this should be moved into the test helper scripts..
	goInstall()

	example_files := map[string][]string{}
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		dir := filepath.Dir(path)
		if dir == "." {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			example_files[dir] = append(example_files[dir], path)
		}
		return err
	})
	for k, files := range example_files {
		if hasTest(files) {
			runTest(t, k)
		}
		if hasNonTest(files) {
			runExample(t, k, files)
		}
	}
}
