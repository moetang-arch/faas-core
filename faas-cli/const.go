package main

import (
	"fmt"
	"os"
)

var (
	faasServiceSrcPath = fmt.Sprint(
		"src",
		string(os.PathSeparator),
		"github.com",
		string(os.PathSeparator),
		"moetang-arch",
		string(os.PathSeparator),
		"faas-core",
		string(os.PathSeparator),
		"faas-service",
	)
)
