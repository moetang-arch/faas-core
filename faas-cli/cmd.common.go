package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
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
			value := v.Path.Value
			if len(value) > 2 {
				value = value[1 : len(value)-1] // omit quotes
				importPaths[value] = struct{}{}
			} else {
				panic(errors.New("import path is empty of file:" + k))
			}
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
LOOP:
	for k := range importPaths {
		_, ok := allowedMap[k]
		if ok {
			continue
		}

		//TODO need a more efficient algorithm
		for pkg := range allowedMap {
			if strings.HasPrefix(k, pkg+string('/')) {
				continue LOOP
			}
		}

		return errors.New("import:" + k + " not allowed")
	}
	return nil
}

func getImports(pkg string) (sourcePackageName string, importPaths map[string]struct{}, err error) {
	defer func() {
		elem := recover()
		if elem != nil {
			switch elem.(type) {
			case error:
				err = elem.(error)
			default:
				err = errors.New("unknown error when processing imports")
			}
		}
	}()
	sourcePackageName, importPaths = checkAndGetImports(pkg)
	return
}

func copyToWithOnlyFiles(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if srcDir == path {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}

		err = copyFile(path, dstDir+string(os.PathSeparator)+info.Name())
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func copyFile(srcFile, dstFile string) error {
	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer f.Close()
	df, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, f)
	if err != nil {
		return err
	}
	return nil
}
