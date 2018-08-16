package main

import (
	"runtime"
	"runtime/debug"
)

func main() {
	// set defaults
	runtime.GOMAXPROCS(1)
	debug.SetMaxStack(4 * 1024 * 1024)
	debug.SetMaxThreads(200)
	//TODO memory/thread/cpu

	service := NewService()
	err := service.Serve()
	if err != nil {
		panic(err)
	}
}
