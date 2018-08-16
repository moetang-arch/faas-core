package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	commands = map[string]func(){
		"run":          run,
		"build":        build,
		"get-core-lib": getCoreLib,
	}
)

func main() {
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if len(os.Args[1]) == 0 || os.Args[1][0] == '-' {
		flag.PrintDefaults()
		//TODO print commands
		os.Exit(2)
	}

	handler, ok := commands[os.Args[1]]
	if !ok {
		fmt.Println(errors.New("unknown command"))
		os.Exit(3)
	}
	handler()
}

func build() {
	//TODO
	fmt.Println("build")
}

func getCoreLib() {
	//TODO get core lib using go get
	fmt.Println("not implement")
}
