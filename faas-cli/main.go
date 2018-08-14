package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	dummy string
)

var (
	commands = map[string]func(){
		"run":          run,
		"compile":      compile,
		"get-core-lib": getCoreLib,
	}
)

func init() {
	flag.StringVar(&dummy, "dummy", "", "do nothing")
}

func main() {
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if len(os.Args[1]) == 0 || os.Args[1][0] == '-' {
		flag.PrintDefaults()
		os.Exit(2)
	}

	err := flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	handler, ok := commands[os.Args[1]]
	if !ok {
		fmt.Println(errors.New("unknown command"))
		os.Exit(4)
	}
	handler()
}

func run() {
	//TODO
	fmt.Println("run")
}

func compile() {
	//TODO
	fmt.Println("compile")
}

func getCoreLib() {
	//TODO get core lib using go get
}
