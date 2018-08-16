package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func checkAndGetImports(srcFolder string) (sourcePackageName string, importPaths map[string]struct{}) {
	fileset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileset, srcFolder, nil, parser.ImportsOnly)
	if err != nil {
		panic(err)
	}

	// package validation
	if len(pkgs) > 2 {
		panic(errors.New("multi-packages defined"))
	}
	if len(pkgs) == 0 {
		panic(errors.New("no package"))
	}

	// package file mapping
	var sourcePkg string
	var sourceFiles = make(map[string]*ast.File)
	var testPkg string
	var testFiles = make(map[string]*ast.File)
	for k, v := range pkgs {
		if strings.HasSuffix(k, "_test") {
			testPkg = k
			for k, v := range v.Files {
				testFiles[k] = v
			}
		} else {
			sourcePkg = k
			for k, v := range v.Files {
				sourceFiles[k] = v
			}
		}
	}

	// file validation
	if len(sourcePkg) == 0 {
		panic(errors.New("no source package"))
	}
	if len(testPkg) > 0 && (sourcePkg+"_test") != testPkg {
		panic(errors.New("illegal source package name or test package name, sourcePkg:" + sourcePkg + " and testPkg:" + testPkg))
	}
	for k := range testFiles {
		k = strings.ToLower(k)
		if !strings.HasSuffix(k, "_test.go") {
			panic(errors.New("illegal test package file name"))
		}
	}

	importPaths = make(map[string]struct{})
	// imports validation
	for k, v := range sourceFiles {
		if strings.HasSuffix(strings.ToLower(k), "_test.go") {
			continue
		}
		for _, v := range v.Imports {
			importPaths[v.Path.Value] = struct{}{}
		}
	}
	sourcePackageName = sourcePkg
	return
}

func checkImportsValidation(importPaths map[string]struct{}, allowedPaths []string) error {
	allowedMap := make(map[string]struct{})
	for _, v := range allowedPaths {
		allowedMap[v] = struct{}{}
	}
	for k, _ := range importPaths {
		_, ok := allowedMap[k]
		if !ok {
			return errors.New("import:" + k + " not allowed")
		}
	}
	return nil
}