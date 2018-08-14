package main

import "runtime"

func main() {
	// set defaults
	runtime.GOMAXPROCS(1)
}
