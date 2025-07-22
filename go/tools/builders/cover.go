// Copyright 2017 The Bazel Authors. All rights reserved.
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func instrumentForCoverage(goenv *env, importPath string, pkgName string, infiles []string, coverVar string, coverMode string, outfiles []string, workDir string) ([]string, error) {
	pkgcfg := workDir + "pkgcfg.txt"
	covoutputs := workDir + "coveroutfiles.txt"
	odir := filepath.Dir(outfiles[0])
	cv := filepath.Join(odir, "covervars.go")
	outfiles = append([]string{cv}, outfiles...)

	pcfg := coverPkgConfig{
		PkgPath:   importPath,
		PkgName:   pkgName,
		Granularity: "perblock",
		OutConfig: pkgcfg,
		Local:     false,
	}
	data, err := json.Marshal(pcfg)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')
	if err := os.WriteFile(pkgcfg, data, 0666); err != nil {
		return nil, err
	}
	var sb strings.Builder
	for i := range outfiles {
		fmt.Fprintf(&sb, "%s\n", outfiles[i])
	}
	if err := os.WriteFile(covoutputs, []byte(sb.String()), 0666); err != nil {
		return nil, err
	}

	goargs := goenv.goTool("cover", "-pkgcfg", pkgcfg, "-var", coverVar, "-mode", coverMode, "-outfilelist", covoutputs)
	goargs = append(goargs, infiles...)
	if err := goenv.runCommand(goargs); err != nil {
		return nil, err
	}
	return outfiles, nil
}

func writeFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0666)
}

// coverPkgConfig matches https://cs.opensource.google/go/go/+/refs/tags/go1.24.4:src/cmd/internal/cov/covcmd/cmddefs.go;l=18
type coverPkgConfig struct {
	// File into which cmd/cover should emit summary info
	// when instrumentation is complete.
	OutConfig string

	// Import path for the package being instrumented.
	PkgPath string

	// Package name.
	PkgName string

	// Instrumentation granularity: one of "perfunc" or "perblock" (default)
	Granularity string

	// Module path for this package (empty if no go.mod in use)
	ModulePath string

	// Local mode indicates we're doing a coverage build or test of a
	// package selected via local import path, e.g. "./..." or
	// "./foo/bar" as opposed to a non-relative import path. See the
	// corresponding field in cmd/go's PackageInternal struct for more
	// info.
	Local bool

	// EmitMetaFile if non-empty is the path to which the cover tool should
	// directly emit a coverage meta-data file for the package, if the
	// package has any functions in it. The go command will pass in a value
	// here if we've been asked to run "go test -cover" on a package that
	// doesn't have any *_test.go files.
	EmitMetaFile string
}
