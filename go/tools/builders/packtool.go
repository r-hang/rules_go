// Copyright 2025 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// packTool builds the pack command from source for the execution platform.
func packTool(args []string) error {
	// Parse arguments.
	args, _, err := expandParamsFiles(args)
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet("GoPackTool", flag.ExitOnError)
	goenv := envFlags(fs)
	var out string
	fs.StringVar(&out, "out", "", "Path to output pack binary")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := goenv.checkFlagsAndSetGoroot(); err != nil {
		return err
	}
	if out == "" {
		return flag.ErrHelp
	}

	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		return fmt.Errorf("GOROOT not set")
	}

	// Create output directory if needed
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}

	// Build cmd/pack using go install for the execution platform
	originalOS, originalARCH := os.Getenv("GOOS"), os.Getenv("GOARCH")
	os.Setenv("GOOS", runtime.GOOS)
	os.Setenv("GOARCH", runtime.GOARCH)
	defer func() {
		os.Setenv("GOOS", originalOS)
		os.Setenv("GOARCH", originalARCH)
	}()

	// Disable modules and set cache
	os.Setenv("GO111MODULE", "off")
	cachePath := filepath.Join(filepath.Dir(out), ".gocache")
	os.Setenv("GOCACHE", cachePath)
	defer os.RemoveAll(cachePath)

	// Disable CGO for building pack tool
	os.Setenv("CGO_ENABLED", "0")

	// Build pack command to temporary location then move to final location
	tempOut := out + ".tmp"
	buildArgs := goenv.goCmd("build", "-o", tempOut, "cmd/pack")
	if err := goenv.runCommand(buildArgs); err != nil {
		return err
	}

	// Move to final location
	return os.Rename(tempOut, out)
}