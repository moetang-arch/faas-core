package main

import (
	"flag"
	"fmt"
	"os"
)

func run() {
	var gopath string
	flag.StringVar(&gopath, "gopath", "", "gopath for faas-api and faas-service source package location")

	err := flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	// working directory
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	packageName, imports, err := getImports(workDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = checkImportsValidation(imports, append(golangPkgAllowedList, userPkgAllowedList...))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// path of faas-core
	// default in gopath
	var faasServiceSourcePath = makeFaasServicePkgPath(getFirstGoPath())

	err = prepareBuildSrc(faasServiceSourcePath, workDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(packageName, workDir, faasServiceSourcePath)
}

func prepareBuildSrc(serviceSourcePath string, curDir string) error {
	dir := curDir + string(os.PathSeparator) + ".build"
	defer func() {
		os.RemoveAll(dir)
	}()
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	dirFile, err := os.Create(dir)
	if err != nil {
		return err
	}
	dirFile.Name()
	//TODO copy dir
	return err
}
