package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	commands = map[string]func(){
		"run":          run,
		"build":        build,
		"get-core-lib": getCoreLib,
	}
	commands2 = map[string]struct {
		Handler     func()
		Description string
	}{
		"run":          {Handler: run, Description: "run faas in local environment"},
		"build":        {Handler: build, Description: "build faas binary file"},
		"get-core-lib": {Handler: getCoreLib, Description: "get core library of faas-core"},
	}
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("available command:")
		// print commands
		for k, v := range commands2 {
			fmt.Println("\t" + k + "\t" + v.Description)
		}
		os.Exit(1)
	}
	if len(os.Args[1]) == 0 || os.Args[1][0] == '-' {
		fmt.Println("unknown command:", os.Args[1])
		fmt.Println("available command:")
		// print commands
		for k, v := range commands2 {
			fmt.Println("\t" + k + "\t" + v.Description)
		}
		os.Exit(2)
	}

	elem, ok := commands2[os.Args[1]]
	if !ok {
		fmt.Println(errors.New("unknown command"))
		os.Exit(3)
	}
	elem.Handler()
}

func build() {
	//TODO
	fmt.Println("build")
}

func getCoreLib() {
	//TODO get core lib using go get
	fmt.Println("not implement")
}
