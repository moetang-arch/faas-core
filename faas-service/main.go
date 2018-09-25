package main

import (
	"runtime"
	"runtime/debug"
)

func main() {
	resourceLimit()

	if err := readingRegistry(); err != nil {
		panic(err)
	}

	service := NewService()
	err := service.Serve()
	if err != nil {
		panic(err)
	}
}

func resourceLimit() {
	// set defaults
	runtime.GOMAXPROCS(1)
	debug.SetMaxStack(4 * 1024 * 1024)
	debug.SetMaxThreads(200)
	//TODO memory/thread/cpu
}
