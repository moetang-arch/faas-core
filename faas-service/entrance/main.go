package main

import (
	"runtime"
	"github.com/moetang-arch/faas-core/faas-service"
)

func main() {
	// set defaults
	runtime.GOMAXPROCS(1)

	service := faas_service.NewService()
}
