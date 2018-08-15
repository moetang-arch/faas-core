package main

import (
	"fmt"
	"os"
	"strings"
)

func getFirstGoPath() string {
	osGoPath := os.Getenv("GOPATH")
	if len(osGoPath) == 0 {
		return ""
	}
	paths := strings.Split(osGoPath, string(os.PathListSeparator))
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func makeFaasServicePkgPath(gopath string) string {
	return fmt.Sprint(gopath, string(os.PathSeparator), faasServiceSrcPath)
}
