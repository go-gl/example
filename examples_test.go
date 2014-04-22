// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const default_wait_time = 1 * time.Second

// Test specific run times
var runtimes = map[string]time.Duration{
	"nehe/03": 5 * time.Second,
}

func Command(args ...string) *exec.Cmd {
	if os.Getenv("USE_VGL") != "" {
		return exec.Command("vglrun", args...)
	}
	return exec.Command(args[0], args[1:]...)
}

func runTest(t *testing.T, path string) {
	println(strings.Repeat("=", 80))
	println("-- subtest: ", path)
	test := Command("go", "test", "-v", "./"+path)
	test.Stdout, test.Stderr = os.Stdout, os.Stderr

	err := test.Run()
	if err != nil {
		t.Fatal("Failed to run 'go test': ", err)
	}
}

func runExample(t *testing.T, path string, files []string) {
	println(strings.Repeat("=", 80))

	get := exec.Command("go", "get", "-d", "-v", "./"+path)
	get.Stdout, get.Stderr = os.Stdout, os.Stderr
	err := get.Run()
	if err != nil {
		t.Fatal("Failed to run go get: ", err)
	}

	bin_name := filepath.Join("bin", path)
	bld := exec.Command("go", "build", "-v", "-o", bin_name, "./"+path)
	bld.Stdout, bld.Stderr = os.Stdout, os.Stderr
	err = bld.Run()
	if err != nil {
		t.Fatal("Failed to run go build: ", err)
	}

	println(strings.Repeat("-", 80))

	cmd := Command(bin_name)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	exited, sleeproutine_error := make(chan bool, 1), make(chan error, 1)
	ran_long_enough := false

	go func() {
		// This go routine waits for the process to exit, or `wait_time`, 
		// whichever comes first. If `wait_time` expires, `cmd` is killed.
		wait_time, ok := runtimes[path]
		if !ok {
			wait_time = default_wait_time
		}
		select {
		case <-exited:
			break
		case <-time.After(wait_time):
			ran_long_enough = true
			err := cmd.Process.Kill()
			if err != nil && err.Error() != "os: process already finished" {
				sleeproutine_error <- err
			}
		}
		sleeproutine_error <- nil
	}()

	err = cmd.Run()
	exited <- true

	if !ran_long_enough && err != nil {
		t.Fatalf("Process died unexpectedly: %q", err)
	}

	// Sync with sleep routine
	if err = <-sleeproutine_error; err != nil {
		t.Fatalf("Unexpected error: %q", err)
	}
}

func goInstall() {
	install := exec.Command("go", "install", "-v")
	install.Stdout, install.Stderr = os.Stdout, os.Stderr
	err := install.Run()
	if err != nil {
		panic(err)
	}
}

// This is needed so that gas can find the go-gl/examples/data/ directory if we
// are running in a clone. Not an ideal solution, better would be for gas to 
// detect the fully qualified package name of the caller. Something to 
// investigate on a rainy day..
func goGetExamples() {
	getex := exec.Command("go", "get", "-d", "github.com/go-gl/examples")
	getex.Stdout, getex.Stderr = os.Stdout, os.Stderr
	err := getex.Run()
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
	goGetExamples()
	goInstall()

	example_files := map[string][]string{}
	// Discover all .go files below this directory
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

	// Run tests and build/run binaries in each directory, if there are any.
	for k, files := range example_files {
		log.Print("=== ", k)
		if hasTest(files) {
			runTest(t, k)
		}
		if hasNonTest(files) {
			runExample(t, k, files)
		}
	}
}
