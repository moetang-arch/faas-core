package main

import (
	"errors"
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

	packageName, imports, err := getImports(makeFaasServicePkgPath(getFirstGoPath()))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = checkImportsValidation(imports, append(golangPkgAllowedList, userPkgAllowedList...))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(packageName)
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
