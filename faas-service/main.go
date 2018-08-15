package main

import (
	"runtime"
)

func main() {
	// set defaults
	runtime.GOMAXPROCS(1)

	service := NewService()
	err := service.Serve()
	if err != nil {
		panic(err)
	}
}
