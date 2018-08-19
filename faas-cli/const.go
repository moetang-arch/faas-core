package main

import (
	"fmt"
	"os"
)

var (
	faasServiceSrcPathForRunCmd = fmt.Sprint(
		"src",
		string(os.PathSeparator),
		"github.com",
		string(os.PathSeparator),
		"moetang-arch",
		string(os.PathSeparator),
		"faas-core",
		string(os.PathSeparator),
		"faas-local",
	)
	faasServiceSrcPathForBuildCmd = fmt.Sprint(
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

	// go 1.11 beta
	golangPkgAllowedList = []string{
		"archive",
		"bufio",
		"builtin",
		"bytes",
		"compress",
		"container",
		"context",
		"crypto",
		"database",
		"debug",
		"encoding",
		"errors",
		"expvar",
		"flag",
		"fmt",
		"go",
		"hash",
		"html",
		"image",
		"index",
		"io",
		"log",
		"math",
		"mime",
		"net",
		"os",
		"path",
		"plugin",
		"reflect",
		"regexp",
		"runtime",
		"sort",
		"strconv",
		"strings",
		"sync",
		"syscall",
		"testing",
		"text",
		"time",
		"unicode",
		"unsafe",
	}

	userPkgAllowedList = []string{
		"github.com/moetang-arch/faas-api",
	}

	generatedSourceTemplate = `
package main

import . "{{.}}"
`
)
